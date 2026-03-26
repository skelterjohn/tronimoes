package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/skelterjohn/tronimoes/tronserv/client"
	"github.com/skelterjohn/tronimoes/tronserv/clog"
)

var (
	addr     = flag.String("addr", "http://localhost:8080", "game server base URL (agent -addr)")
	agentExe = flag.String("agent", "agent", "agent executable: name on PATH or full path")
	logDir   = flag.String("logdir", "arena_logs", "directory for run logs; a subdir YYYYMMDD_HHMMSS_<gamecode> is created under it")
)

// safeFilename returns a filesystem-safe basename (same idea as gibbs_planner/evaluate).
func safeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, name)
}

func randomGameCodeAZ6() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func expandCode(ctx context.Context, code string) string {
	if len(code) > 6 {
		return code
	}

	tc := &client.TronimoesClient{
		TronservAddr: *addr,
		Client:       http.DefaultClient,
		Code:         code,
	}
	g, err := tc.GetGame(ctx, 0)
	if err != nil {
		return code
	}
	clog.Info(ctx, "Expanded code", "code", code, "new code", g.Code)
	return g.Code
}

func main() {
	ctx := context.Background()

	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("usage: arena [options] config.yaml")
		flag.PrintDefaults()
		exitCode = 2
		return
	}

	ac, err := LoadArenaConfig(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		exitCode = 1
		return
	}

	var started []*exec.Cmd

	gameCode := randomGameCodeAZ6()
	fmt.Println("game code", gameCode)

	runDir := filepath.Join(*logDir, time.Now().Format("20060102_150405")+"_"+gameCode)
	if err := os.MkdirAll(runDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir run log dir %s: %v\n", runDir, err)
		exitCode = 1
		return
	}
	runDir, err = filepath.Abs(runDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abs run log dir %s: %v\n", runDir, err)
		exitCode = 1
		return
	}
	fmt.Println("agent stdout logs written to", runDir)

	fmt.Println("spawning agents, connecting to", *addr)

	for _, pcfg := range ac.Players {
		name := pcfg.Name
		logPath := filepath.Join(runDir, safeFilename(name))
		logF, err := os.Create(logPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create agent log file %s: %v\n", logPath, err)
			exitCode = 1
			return
		}
		defer logF.Close()

		cfgBuf := &bytes.Buffer{}
		if err := json.NewEncoder(cfgBuf).Encode(pcfg.Config); err != nil {
			fmt.Fprintf(os.Stderr, "encode agent config: %v\n", err)
			exitCode = 1
			return
		}

		fmt.Printf("%s config: %s\n", name, cfgBuf.String())

		cmd := exec.Command(*agentExe,
			"--addr", *addr,
			"--name", name,
			"--code", gameCode,
			"--which", "gibbs",
			"--no-react",
			"--dev",
			"--ready-with", fmt.Sprintf("%d", len(ac.Players)),
			"--min-move-time", "0",
			"--config", "-",
		)
		cmd.Stdout = logF
		cmd.Stderr = os.Stderr
		cmd.Stdin = cfgBuf
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "start agent %s: %v\n", name, err)
			exitCode = 1
			return
		}
		defer cmd.Process.Kill()
		started = append(started, cmd)
		for len(gameCode) == 6 {
			time.Sleep(1 * time.Second)
			gameCode = expandCode(ctx, gameCode)
			if len(gameCode) != 6 {
				fmt.Println("expanded game code to", gameCode)
			}
		}
	}

	fmt.Println("game board", fmt.Sprintf("http://localhost:3000/gameboard/%s", gameCode))

	var wg sync.WaitGroup
	var runErr error
	var mu sync.Mutex
	for _, cmd := range started {
		wg.Add(1)
		go func(c *exec.Cmd) {
			defer wg.Done()
			if err := c.Wait(); err != nil {
				mu.Lock()
				if runErr == nil {
					runErr = err
				}
				mu.Unlock()
			}
		}(cmd)
	}
	wg.Wait()

	if runErr != nil {
		fmt.Fprintf(os.Stderr, "agent exited with error: %v\n", runErr)
		exitCode = 1
		return
	}
}
