package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/skelterjohn/tronimoes/tronserv/clog"
)

func prefixLines(prefix string, r io.Reader, w io.Writer, mu *sync.Mutex) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		mu.Lock()
		fmt.Fprintf(w, "%s: %s\n", prefix, line)
		mu.Unlock()
	}
}

func replicate(ctx context.Context, prefix string, args []string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		clog.Info(ctx, "Could not create stdout pipe", "error", err.Error(), "prefix", prefix)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		clog.Info(ctx, "Could not create stderr pipe", "error", err.Error(), "prefix", prefix)
		return
	}
	if err := cmd.Start(); err != nil {
		clog.Info(ctx, "Could not start command", "error", err.Error(), "prefix", prefix)
		return
	}

	var mu sync.Mutex
	var streamWg sync.WaitGroup
	streamWg.Add(2)
	go func() {
		defer streamWg.Done()
		prefixLines(prefix, stdout, os.Stdout, &mu)
	}()
	go func() {
		defer streamWg.Done()
		prefixLines(prefix, stderr, os.Stderr, &mu)
	}()

	streamWg.Wait()
	if err := cmd.Wait(); err != nil {
		clog.Info(ctx, "Command exited with error", "error", err.Error(), "prefix", prefix)
	}
}

func main() {
	ctx := context.Background()
	ctx = clog.WithStructuredOutput(ctx, os.Stdout)
	ctx = clog.WithSeverities(ctx, "info")

	count, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		clog.Fatal(ctx, "Could not parse count", err)
	}
	args := os.Args[2:]

	clog.Info(ctx, "Starting replicants", "count", fmt.Sprint(count))

	var wg sync.WaitGroup
	for i := int64(0); i < count; i++ {
		wg.Add(1)
		go replicate(ctx, fmt.Sprintf("replicant-%d", i), args, &wg)
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	clog.Info(ctx, "All replicants have finished")
}
