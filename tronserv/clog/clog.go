package clog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func Info(ctx context.Context, message string, addTags ...string) {
	if !severityEnabled(ctx, infoKey) {
		return
	}
	Log(ctx, "INFO", message, addTags...)
}

func Error(ctx context.Context, message string, addTags ...string) {
	if !severityEnabled(ctx, errorKey) {
		return
	}
	Log(ctx, "ERROR", message, addTags...)
}

func Debug(ctx context.Context, message string, addTags ...string) {
	if !severityEnabled(ctx, debugKey) {
		return
	}
	Log(ctx, "DEBUG", message, addTags...)
}

func Log(ctx context.Context, severity, message string, addTags ...string) {
	existingTags, ok := ctx.Value(keywordsKey).(map[string]string)
	if !ok {
		existingTags = make(map[string]string)
	}
	for i, key := range addTags {
		value := ""
		if len(addTags) > i {
			value = addTags[i+1]
		}
		existingTags[key] = value
	}
	fl := fileline()
	if durationSince, ok := ctx.Value(durationsKey).(time.Time); ok {
		duration := time.Since(durationSince)
		existingTags["duration"] = fmt.Sprintf("%6dms", duration.Milliseconds())
		existingTags["since"] = durationSince.Format(time.RFC3339)
	}
	if textOutput, ok := ctx.Value(textOutputKey).(io.Writer); ok {
		strb := strings.Builder{}
		strb.WriteString(fl)
		strb.WriteString("\t")
		strb.WriteString("| ")
		for key, value := range existingTags {
			strb.WriteString(key)
			strb.WriteString("=")
			strb.WriteString(value)
			strb.WriteString(" ")
		}
		strb.WriteString("|\t")
		strb.WriteString(message)
		fmt.Fprintln(textOutput, strb.String())
	}
	if structuredOutput, ok := ctx.Value(structuredOutputKey).(io.Writer); ok {
		if err := json.NewEncoder(structuredOutput).Encode(map[string]any{
			"severity": severity,
			"message":  message,
			"tags":     existingTags,
			"fileline": fileline(),
		}); err != nil {
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
