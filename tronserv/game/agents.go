package game

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type AgentSpawner interface {
	NewAgent(ctx context.Context, which string, code string) error
}

type LocalAgentSpawner struct {
}

func (s LocalAgentSpawner) NewAgent(ctx context.Context, which string, code string) error {
	log.Printf("Spawning %s agent for game %s", which, code)
	exeDir := ""
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
	}
	// Use absolute path so Windows allows running it (avoids "cannot run executable found relative to current directory").
	agentExe := filepath.Join(exeDir, "tronagent.exe")
	// Use Background so the agent is not killed when the HTTP request context is cancelled.
	runCtx := context.WithoutCancel(ctx)
	cmd := exec.CommandContext(runCtx,
		agentExe,
		"--which", which,
		"--code", code,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = exeDir
	if exeDir != "" {
		log.Printf("Running agent in %s", exeDir)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Agent %q for %q exited: %v", which, code, err)
		}
	}()
	return nil
}
