package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/skelterjohn/tronimoes/tronserv/clog"
)

func replicate(ctx context.Context, args []string) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		clog.Error(ctx, "Could not run command", err)
		return
	}
}

func main() {
	dev := false
	for _, arg := range os.Args[1:] {
		if arg == "--dev" {
			dev = true
			break
		}
	}
	ctx := context.Background()
	if dev {
		ctx = clog.WithTextOutput(ctx, os.Stdout)
	} else {
		ctx = clog.WithStructuredOutput(ctx, os.Stdout)
	}
	ctx = clog.WithSeverities(ctx, clog.INFO)

	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		clog.Fatal(ctx, "Could not parse count", err)
	}
	args := os.Args[2:]

	clog.Info(ctx, "Starting replicants", "count", fmt.Sprint(count))

	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			replicate(ctx, args)
		}()
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	clog.Info(ctx, "All replicants have finished")
}
