package gibbs_planner

import (
	"context"
	"fmt"
	"io"
	"os"
)

type logKey struct{}

// WithLogBuffer returns a context that stores w as the log destination.
// Log(ctx, ...) will write to w when given a ctx derived from this.
func WithLogBuffer(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, logKey{}, w)
}

// Log writes to the writer in ctx (if set via WithLogBuffer), otherwise to stdout.
func Log(ctx context.Context, format string, args ...interface{}) {
	w := ctx.Value(logKey{})
	if w != nil {
		if bw, ok := w.(io.Writer); ok {
			fmt.Fprintf(bw, format+"\n", args...)
			return
		}
	}
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}
