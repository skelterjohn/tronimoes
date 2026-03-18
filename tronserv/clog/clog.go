package clog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"cloud.google.com/go/logging"
)

type keywordsKeyType string
type textOutputKeyType string
type structuredOutputKeyType string
type cloudLoggingKeyType string
type infoKeyType string
type errorKeyType string
type debugKeyType string
type durationsKeyType string

var keywordsKey keywordsKeyType = "keywords"
var textOutputKey textOutputKeyType = "textOutput"
var structuredOutputKey structuredOutputKeyType = "structuredOutput"
var cloudLoggingKey cloudLoggingKeyType = "cloudLogging"
var infoKey infoKeyType = "info"
var errorKey errorKeyType = "error"
var debugKey debugKeyType = "debug"
var durationsKey durationsKeyType = "durations"

func WithTextOutput(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, textOutputKey, w)
}

func WithStructuredOutput(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, structuredOutputKey, w)
}

// cloudLoggingValue holds the Cloud Logging client and logger created by WithCloudLoggingOutput.
type cloudLoggingValue struct {
	Logger *logging.Logger
	Client *logging.Client
}

// WithCloudLoggingOutput creates a Cloud Logging client and logger using logID (project ID
// is auto-detected on GCP) and stores them in the context. When set, _log writes entries
// to the Google Cloud Logging API. Call CloseCloudLogging before exit to flush buffered entries.
func WithCloudLoggingOutput(ctx context.Context, logID string) context.Context {
	if logID == "" {
		logID = "app"
	}
	client, err := logging.NewClient(ctx, logging.DetectProjectID)
	if err != nil {
		log.Fatalf("Could not create Cloud Logging client: %v", err)
	}
	return context.WithValue(ctx, cloudLoggingKey, &cloudLoggingValue{
		Logger: client.Logger(logID),
		Client: client,
	})
}

// CloseCloudLogging closes the Cloud Logging client stored in the context, flushing
// any buffered entries. Call before process exit when using WithCloudLoggingOutput.
func CloseCloudLogging(ctx context.Context) {
	if v, ok := ctx.Value(cloudLoggingKey).(*cloudLoggingValue); ok && v != nil && v.Client != nil {
		if err := v.Client.Close(); err != nil {
			log.Printf("Could not close Cloud Logging client: %v", err)
		}
	}
}

func WithSeverities(ctx context.Context, severities ...string) context.Context {
	for _, severity := range severities {
		switch severity {
		case "info":
			ctx = context.WithValue(ctx, infoKey, true)
		case "error":
			ctx = context.WithValue(ctx, errorKey, true)
		case "debug":
			ctx = context.WithValue(ctx, debugKey, true)
		}
	}
	return ctx
}

func WithKeyword(ctx context.Context, key, value string) context.Context {
	existingTags, ok := ctx.Value(keywordsKey).(map[string]string)
	if !ok {
		existingTags = make(map[string]string)
	}
	existingTags[key] = value
	return context.WithValue(ctx, keywordsKey, existingTags)
}

func WithDurationSince(ctx context.Context, start time.Time) context.Context {
	return context.WithValue(ctx, durationsKey, start)
}

func severityEnabled(ctx context.Context, key any) bool {
	v, _ := ctx.Value(key).(bool)
	return v
}

func Info(ctx context.Context, message string, addTags ...any) {
	if !severityEnabled(ctx, infoKey) {
		return
	}
	_log(ctx, "INFO", message, tagsFromList(addTags...))
}

func Error(ctx context.Context, message string, err error, addTags ...any) {
	if !severityEnabled(ctx, errorKey) {
		return
	}
	tags := tagsFromList(addTags...)
	if err != nil {
		tags["error"] = err.Error()
	}
	_log(ctx, "ERROR", message, tags)
}

func Debug(ctx context.Context, message string, addTags ...any) {
	if !severityEnabled(ctx, debugKey) {
		return
	}
	_log(ctx, "DEBUG", message, tagsFromList(addTags...))
}

func Fatal(ctx context.Context, message string, err error, addTags ...any) {
	tags := tagsFromList(addTags...)
	if err != nil {
		tags["error"] = err.Error()
	}
	_log(ctx, "FATAL", message, tags)
	os.Exit(1)
}

func valueString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	}
	return fmt.Sprint(value)
}

func tagsFromList(addTags ...any) map[string]string {
	tags := make(map[string]string)
	for i, key := range addTags {
		if i%2 == 1 {
			continue
		}
		keyStr, ok := key.(string)
		if !ok {
			log.Printf("tag %v is not a string", key)
			continue
		}
		var value any
		if len(addTags) > i+1 {
			value = addTags[i+1]
		} else {
			value = "UNKNOWN"
		}
		tags[keyStr] = valueString(value)
	}
	return tags
}

func Log(ctx context.Context, severity, message string, addTags map[string]any) {
	tags := make(map[string]string)
	for key, value := range addTags {
		tags[key] = valueString(value)
	}
	_log(ctx, severity, message, tags)
}

func _log(ctx context.Context, severity, message string, addTags map[string]string) {
	existingTags, _ := ctx.Value(keywordsKey).(map[string]string)
	tagsForMessage := make(map[string]string)
	maps.Copy(tagsForMessage, existingTags)
	maps.Copy(tagsForMessage, addTags)
	fl := fileline()
	if durationSince, ok := ctx.Value(durationsKey).(time.Time); ok {
		duration := time.Since(durationSince)
		tagsForMessage["duration"] = fmt.Sprintf("%6dms", duration.Milliseconds())
		tagsForMessage["since"] = durationSince.Format(time.RFC3339)
	}
	if textOutput, ok := ctx.Value(textOutputKey).(io.Writer); ok {
		strb := strings.Builder{}
		strb.WriteString(fl)
		strb.WriteString("\t")
		strb.WriteString(message)
		strb.WriteString("\t| ")
		// Include both tags already on the context (keywords) and tags passed to Log().
		for key, value := range tagsForMessage {
			strb.WriteString(" ")
			strb.WriteString(key)
			strb.WriteString("=")
			strb.WriteString(value)
		}
		fmt.Fprintln(textOutput, strb.String())
	}
	if structuredOutput, ok := ctx.Value(structuredOutputKey).(io.Writer); ok {
		allTags := map[string]any{
			"severity": severity,
			"message":  message,
			"fileline": fileline(),
		}
		// Include both tags already on the context (keywords) and tags passed to Log().
		for key, value := range tagsForMessage {
			allTags[key] = value
		}
		if err := json.NewEncoder(structuredOutput).Encode(allTags); err != nil {
			log.Printf("could not marshal to structured output: %v", err)
		}
	}
	if v, ok := ctx.Value(cloudLoggingKey).(*cloudLoggingValue); ok && v != nil && v.Logger != nil {
		allTags := map[string]string{
			"severity": severity,
			"message":  message,
			"fileline": fileline(),
		}
		// Include both tags already on the context (keywords) and tags passed to Log().
		maps.Copy(allTags, tagsForMessage)
		v.Logger.Log(logging.Entry{
			Severity: severityFromString(severity),
			Labels:   allTags,
			Payload:  message,
		})
	}
}

func severityFromString(s string) logging.Severity {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return logging.Debug
	case "INFO":
		return logging.Info
	case "NOTICE":
		return logging.Notice
	case "WARNING":
		return logging.Warning
	case "ERROR":
		return logging.Error
	case "CRITICAL":
		return logging.Critical
	case "ALERT":
		return logging.Alert
	case "EMERGENCY":
		return logging.Emergency
	case "FATAL":
		return logging.Critical
	default:
		return logging.Default
	}
}

func Print(ctx context.Context, format string, items ...any) {
	msg := fmt.Sprintf(format, items...)
	if tag, ok := ctx.Value(keywordsKey).(string); ok {
		msg = fmt.Sprintf("[%s] %s", tag, msg)
	}
	log.Printf("%s\t| %s", fileline(), msg)
}

func fileline() string {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		file = filepath.Base(file)
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}

func JSON(v any) string {
	d, err := json.Marshal(v)
	if err != nil {
		log.Printf("could not marshal from %s: %v", fileline(), err)
	}
	return string(d)
}
