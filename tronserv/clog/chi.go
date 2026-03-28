package clog

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// chiLogFormatter implements chi's LogFormatter and logs each request via clog.
type chiLogFormatter struct{}

func (chiLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &chiLogEntry{request: r}
}

type chiLogEntry struct {
	request *http.Request
}

func (e *chiLogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	Log(e.request.Context(), "INFO", "serving request", map[string]any{
		"method":   e.request.Method,
		"path":     e.request.URL.RequestURI(),
		"status":   status,
		"bytes":    bytes,
		"duration": elapsed,
	})
}

func (e *chiLogEntry) Panic(v interface{}, _ []byte) {
	Log(e.request.Context(), "ERROR", "request panic",
		map[string]any{
			"panic": v,
		})
}

// ChiLogFormatter returns a chi LogFormatter that logs requests using clog.Log
// with "INFO" or "ERROR" as appropriate. ChiLogger wraps this and injects output
// so you typically use ChiLogger instead.
func ChiLogFormatter() middleware.LogFormatter {
	return chiLogFormatter{}
}

// ChiLogger returns chi middleware that injects clog into the request context
// (structured output to os.Stdout, severities INFO and ERROR for handlers)
// and logs each request via ChiLogFormatter using Log() with "INFO" or "ERROR".
// Use as r.Use(clog.ChiLogger()).
func ChiLogger(ctx context.Context, withServerLogs bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		inject := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Merge clog into the request context; replacing r.Context() drops chi's
				// route context and panics inside chi.Mux.routeHTTP.
				next.ServeHTTP(w, r.WithContext(mergeClogContext(r.Context(), ctx)))
			})
		}
		if withServerLogs {
			return inject(middleware.RequestLogger(ChiLogFormatter())(next))
		} else {
			return inject(next)
		}
	}
}
