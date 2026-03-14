package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
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
		log.Printf("[%s] Could not create stdout pipe: %v", prefix, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("[%s] Could not create stderr pipe: %v", prefix, err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("[%s] Could not start command: %v", prefix, err)
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
		log.Printf("[%s] Command exited with error: %v", prefix, err)
	}
}

func main() {
	ctx := context.Background()

	count, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("Could not parse count: %v", err)
	}
	args := os.Args[2:]

	log.Printf("Starting %d replicants via %v", count, args)

	var wg sync.WaitGroup
	for i := int64(0); i < count; i++ {
		wg.Add(1)
		go replicate(ctx, fmt.Sprintf("replicant-%d", i), args, &wg)
	}
	wg.Wait()
	log.Println("All replicants have finished")
}
