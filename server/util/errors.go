package util

import (
	"fmt"

	"google.golang.org/grpc/status"
)

func Annotate(err error, format string, items ...interface{}) error {
	upstream := err.Error()
	if s, ok := status.FromError(err); ok {
		upstream = s.Message()
	}
	msg := fmt.Sprintf(format, items...) + ": " + upstream
	return status.Error(status.Code(err), msg)
}
