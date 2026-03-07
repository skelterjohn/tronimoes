package gibbs_planner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

func loadGameFromTestdata(t *testing.T, label string) *game.Game {
	path := filepath.Join("testdata", label+".json")
	encoded, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", path, err)
	}
	var g game.Game
	if err := json.Unmarshal(encoded, &g); err != nil {
		t.Fatalf("could not unmarshal: %v", err)
	}
	return &g
}

func TestOneshot(t *testing.T) {
	g := loadGameFromTestdata(t, "oneshot")
	ctx := t.Context()

	// Agent plays as the current player.
	currentPlayer := g.Players[g.Turn]
	gp := &GibbsPlanner{
		Name:               currentPlayer.Name,
		MaxInferenceTime:   1 * time.Second,
		MaxSimulationTime:  1 * time.Second,
		MaxSimulationDepth: 4,
		ValueDecay:         0.9,
	}

	// Previous game has no rounds so Update runs createInitialGuesses.
	previousGame := &game.Game{
		Players: g.Players,
		MaxPips: g.MaxPips,
	}

	gp.Update(ctx, previousGame, g)

	move := gp.GetMove(ctx, g, currentPlayer)
	if move.LaidTile == nil {
		t.Fatalf("Did not play a one-shot: %s", move)
	}

	move.LaidTile.PlayerName = currentPlayer.Name
	if err := g.LayTile(ctx, currentPlayer.Name, move.LaidTile); err != nil {
		t.Fatalf("could not lay tile: %v", err)
	}

	if r := g.CurrentRound(ctx); r != nil {
		t.Fatalf("Round is not done after move: %s", move)
	}

}
