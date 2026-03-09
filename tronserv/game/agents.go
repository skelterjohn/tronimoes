package game

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type AgentSpawner interface {
	NewAgent(ctx context.Context, which string, code string) error
}

type LocalAgentSpawner struct {
}

func (s LocalAgentSpawner) NewAgent(ctx context.Context, which string, code string) error {
	log.Printf("Spawning agent %q for %q", which, code)
	cmd := exec.CommandContext(ctx,
		"go", "run", "github.com/skelterjohn/tronimoes/tronserv/agent",
		"--which", which,
		"--code", code,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}
	return nil
}
