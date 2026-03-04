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
	path := filepath.Join("../../game/testdata", label+".json")
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

func TestGibbsUpdateAndGetMove(t *testing.T) {
	g := loadGameFromTestdata(t, "basic_report")
	ctx := t.Context()

	// Agent plays as the current player.
	currentPlayer := g.Players[g.Turn]
	gp := &GibbsPlanner{
		Name:               currentPlayer.Name,
		MaxInferenceTime:   1 * time.Second,
		MaxSimulationTime:  1 * time.Second,
		MaxSimulationDepth: 10,
	}

	// Previous game has no rounds so Update runs createInitialGuesses.
	previousGame := &game.Game{
		Players: g.Players,
		MaxPips: g.MaxPips,
	}

	gp.Update(ctx, previousGame, g)

	move := gp.GetMove(ctx, g, currentPlayer)
	t.Errorf("move: %v", move)
	if move.Draw && !currentPlayer.JustDrew {
		t.Logf("GetMove returned Draw (no legal lay/spacer)")
	}
	t.Logf("move: Draw=%v Pass=%v LaidTile=%v Spacer=%v",
		move.Draw, move.Pass, move.LaidTile, move.Spacer)
}

func TestShortSimulation(t *testing.T) {
	ctx := t.Context()
	g := game.NewGame(ctx, "AAAAAA")
	if err := g.AddPlayer(ctx, &game.Player{Name: "A"}); err != nil {
		t.Fatalf("error adding player A: %v", err)
	}
	if err := g.AddPlayer(ctx, &game.Player{Name: "B"}); err != nil {
		t.Fatalf("error adding player B: %v", err)
	}
	if err := g.Start(ctx, "A"); err != nil {
		t.Fatalf("error starting game: %v", err)
	}
	if err := g.Start(ctx, "B"); err != nil {
		t.Fatalf("error starting game: %v", err)
	}

	gp := &GibbsPlanner{
		Name:                  "A",
		MaxInferenceTime:      1 * time.Second,
		MaxSimulationTime:     1 * time.Second,
		MaxSimulationDepth:    10,
		MaxSimulationsPerMove: 0,
	}

	gp.Update(ctx, nil, g)

	move := gp.GetMove(ctx, g, g.Players[g.Turn])
	t.Logf("move: Draw=%v Pass=%v LaidTile=%v Spacer=%v",
		move.Draw, move.Pass, move.LaidTile, move.Spacer)
	t.Error(move.String())
}
