package clog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type keywordsKeyType string
type textOutputKeyType string
type structuredOutputKeyType string
type infoKeyType string
type errorKeyType string
type debugKeyType string
type durationsKeyType string

var keywordsKey keywordsKeyType = "keywords"
var textOutputKey textOutputKeyType = "textOutput"
var structuredOutputKey structuredOutputKeyType = "structuredOutput"
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
	switch value.(type) {
	case string:
		return value.(string)
	case int:
		return fmt.Sprintf("%d", value.(int))
	case bool:
		return fmt.Sprintf("%t", value.(bool))
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
	existingTags, ok := ctx.Value(keywordsKey).(map[string]string)
	if !ok {
		existingTags = make(map[string]string)
	}
	fl := fileline()
	for key, value := range addTags {
		existingTags[key] = value
	}
	if durationSince, ok := ctx.Value(durationsKey).(time.Time); ok {
		duration := time.Since(durationSince)
		existingTags["duration"] = fmt.Sprintf("%6dms", duration.Milliseconds())
		existingTags["since"] = durationSince.Format(time.RFC3339)
	}
	if textOutput, ok := ctx.Value(textOutputKey).(io.Writer); ok {
		strb := strings.Builder{}
		strb.WriteString(fl)
		strb.WriteString("\t")
		strb.WriteString(message)
		strb.WriteString("\t| ")
		for key, value := range existingTags {
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
		for key, value := range existingTags {
			allTags[key] = value
		}
		allTags["fileline"] = fileline()
		if err := json.NewEncoder(structuredOutput).Encode(allTags); err != nil {
			log.Printf("could not marshal to structured output: %v", err)
		}
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
