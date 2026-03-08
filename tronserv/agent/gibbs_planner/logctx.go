package gibbs_planner

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type logKey struct{}

// logOpts is stored in context under logKey. It holds the log writer and optional test start time.
type logOpts struct {
	w     io.Writer
	start time.Time
}

// WithLogBuffer returns a context that stores w as the log destination.
// Log(ctx, ...) will write to w when given a ctx derived from this.
func WithLogBuffer(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, logKey{}, &logOpts{w: w})
}

// WithLogStart stores the test start time in the log opts in ctx (from a prior WithLogBuffer).
// If set, Log() prefixes each line with [Nms] (milliseconds since start).
func WithLogStart(ctx context.Context, start time.Time) context.Context {
	val := ctx.Value(logKey{})
	if opts, _ := val.(*logOpts); opts != nil {
		opts = &logOpts{w: opts.w, start: start}
		return context.WithValue(ctx, logKey{}, opts)
	}
	return context.WithValue(ctx, logKey{}, &logOpts{start: start})
}

// Log writes to the writer in ctx (if set via WithLogBuffer), otherwise to stdout.
// If a start time was set via WithLogStart, the line is prefixed with [Nms] where N is milliseconds since start.
func Log(ctx context.Context, format string, args ...interface{}) {
	val := ctx.Value(logKey{})
	opts, _ := val.(*logOpts)
	var line string
	if opts != nil && !opts.start.IsZero() {
		ms := time.Since(opts.start).Milliseconds()
		line = fmt.Sprintf("[%6dms] "+format+"\n", append([]interface{}{ms}, args...)...)
	} else {
		line = fmt.Sprintf(format+"\n", args...)
	}
	if opts != nil && opts.w != nil {
		fmt.Fprint(opts.w, line)
		return
	}
	fmt.Fprint(os.Stdout, line)
}
