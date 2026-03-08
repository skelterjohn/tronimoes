package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/skelterjohn/tronimoes/tronserv/agent/gibbs_planner"
	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

// SuccessFunc is called after loading a game, running the agent, and applying its move.
// It returns true if the outcome is considered success, and an optional message.
type SuccessFunc func(g *game.Game, move types.Move) (success bool, message string)

// TestCase pairs a start state (testdata JSON label) with a success predicate.
type TestCase struct {
	Name    string
	Label   string // filename in testdata/ without .json
	Success SuccessFunc
}

func loadGame(testdataDir, label string) (*game.Game, error) {
	path := filepath.Join(testdataDir, label+".json")
	encoded, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var g game.Game
	if err := json.Unmarshal(encoded, &g); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &g, nil
}

func runCase(ctx context.Context, testdataDir string, tc TestCase, gp *gibbs_planner.GibbsPlanner) (success bool, message string) {
	g, err := loadGame(testdataDir, tc.Label)
	if err != nil {
		return false, err.Error()
	}

	currentPlayer := g.Players[g.Turn]
	gp.Name = currentPlayer.Name
	previousGame := &game.Game{
		Players: g.Players,
		MaxPips: g.MaxPips,
	}

	gp.Update(ctx, previousGame, g)
	move := gp.GetMove(ctx, g, currentPlayer)

	if move.LaidTile != nil {
		move.LaidTile.PlayerName = currentPlayer.Name
		if err := g.LayTile(ctx, currentPlayer.Name, move.LaidTile); err != nil {
			return false, fmt.Sprintf("LayTile: %v", err)
		}
	} else if move.Draw {
		if !g.DrawTile(ctx, currentPlayer.Name) {
			return false, "DrawTile failed"
		}
	}
	// Pass is not applied here for simplicity; add if needed.

	return tc.Success(g, move)
}

func listNames(tests []TestCase) string {
	names := make([]string, len(tests))
	for i := range tests {
		names[i] = tests[i].Name
	}
	return strings.Join(names, ", ")
}

// safeFilename returns a filesystem-safe name for the test case.
func safeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, name)
}

func main() {
	testsFlag := flag.String("tests", "", "comma-separated list of test names (e.g. Oneshot,NoSelfKill); empty runs all")
	countFlag := flag.Int("count", 1, "run each test this many times")
	concurrencyFlag := flag.Int("concurrency", 1, "run this many tests at a time")
	flag.Parse()

	testdataDir := "testdata"

	allTests := []TestCase{
		{
			Name:  "Oneshot",
			Label: "oneshot",
			Success: func(g *game.Game, move types.Move) (bool, string) {
				if move.LaidTile == nil {
					return false, "expected a lay (one-shot win), got no tile"
				}
				if r := g.CurrentRound(context.Background()); r != nil {
					return false, "round should be done after winning move"
				}
				return true, ""
			},
		},
		{
			Name:  "NoSelfKill",
			Label: "noselfkill",
			Success: func(g *game.Game, move types.Move) (bool, string) {
				if move.LaidTile == nil && !move.Draw {
					return false, "expected a lay or draw, got neither"
				}
				if r := g.CurrentRound(context.Background()); r == nil {
					return false, "round should not be done (player must not kill own line)"
				}
				return true, ""
			},
		},
	}

	// Filter by -tests if set.
	tests := allTests
	if *testsFlag != "" {
		want := make(map[string]bool)
		for _, name := range strings.Split(*testsFlag, ",") {
			want[strings.TrimSpace(name)] = true
		}
		tests = nil
		for _, tc := range allTests {
			if want[tc.Name] {
				tests = append(tests, tc)
			}
		}
		if len(tests) == 0 {
			fmt.Fprintf(os.Stderr, "no tests matched -tests=%q (available: %s)\n", *testsFlag, listNames(allTests))
			os.Exit(1)
		}
	}

	ctx := context.Background()
	count := *countFlag
	concurrency := *concurrencyFlag
	if count < 1 {
		count = 1
	}
	if concurrency < 1 {
		concurrency = 1
	}

	logDir, err := os.MkdirTemp("", "gibbs-eval-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create temp dir for logs: %v\n", err)
		os.Exit(1)
	}
	logDir, _ = filepath.Abs(logDir)
	fmt.Fprintf(os.Stderr, "Logs will be written to: %s\n", logDir)

	// Build jobs: (test case, run index). Run number is zero-padded to runWidth in filenames.
	runWidth := len(fmt.Sprintf("%d", count-1))
	type job struct {
		tc       TestCase
		run      int
		runWidth int
	}
	var jobs []job
	for _, tc := range tests {
		for run := 0; run < count; run++ {
			jobs = append(jobs, job{tc: tc, run: run, runWidth: runWidth})
		}
	}

	// Run with limited concurrency. Each job uses its own planner, log buffer, and loaded game.
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	type result struct {
		name    string
		run     int
		success bool
		msg     string
	}
	results := make(chan result, len(jobs))
	for _, j := range jobs {
		wg.Add(1)
		sem <- struct{}{}
		go func(j job) {
			defer wg.Done()
			defer func() { <-sem }()
			var logBuf bytes.Buffer
			startTime := time.Now()
			runCtx := gibbs_planner.WithLogBuffer(ctx, &logBuf)
			runCtx = gibbs_planner.WithLogStart(runCtx, startTime)
			gp := &gibbs_planner.GibbsPlanner{
				MaxInferenceTime:   1 * time.Second,
				MaxSimulationTime:  1 * time.Second,
				MaxSimulationDepth: 4,
				ValueDecay:         0.9,
			}
			ok, msg := runCase(runCtx, testdataDir, j.tc, gp)
			verdict := "OK"
			if !ok {
				verdict = "FAIL"
				fmt.Fprintf(&logBuf, "\n--- result: %s ---\n", msg)
			}
			fname := fmt.Sprintf("%s_%0*d_%s.log", safeFilename(j.tc.Name), j.runWidth, j.run, verdict)
			path := filepath.Join(logDir, fname)
			_ = os.WriteFile(path, logBuf.Bytes(), 0644)
			results <- result{name: j.tc.Name, run: j.run, success: ok, msg: msg}
		}(j)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate and print.
	byName := make(map[string]struct {
		pass, fail int
		lastFail   string
	})
	for r := range results {
		ent := byName[r.name]
		if r.success {
			ent.pass++
		} else {
			ent.fail++
			ent.lastFail = r.msg
		}
		byName[r.name] = ent
	}
	for _, tc := range tests {
		ent := byName[tc.Name]
		total := ent.pass + ent.fail
		if ent.fail == 0 {
			fmt.Printf("PASS %s (%d/%d)\n", tc.Name, ent.pass, total)
		} else {
			fmt.Printf("FAIL %s (%d/%d passed): %s\n", tc.Name, ent.pass, total, ent.lastFail)
		}
	}
	fmt.Fprintf(os.Stderr, "Logs written to: %s\n", logDir)
}
