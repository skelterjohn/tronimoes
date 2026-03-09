package gibbs_planner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

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
	t.Parallel()
	g := loadGameFromTestdata(t, "oneshot")
	ctx := t.Context()

	// Agent plays as the current player.
	currentPlayer := g.Players[g.Turn]
	gp := &GibbsPlanner{
		Name: currentPlayer.Name,
	}
	gp.SetDefaults()

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

func TestNoSelfKill(t *testing.T) {
	t.Parallel()
	g := loadGameFromTestdata(t, "noselfkill")
	ctx := t.Context()

	currentPlayer := g.Players[g.Turn]
	gp := &GibbsPlanner{
		Name: currentPlayer.Name,
	}
	gp.SetDefaults()

	previousGame := &game.Game{
		Players: g.Players,
		MaxPips: g.MaxPips,
	}

	gp.Update(ctx, previousGame, g)

	move := gp.GetMove(ctx, g, currentPlayer)
	if move.LaidTile == nil && !move.Draw {
		t.Fatalf("Did not play a tile or draw a tile: %s", move)
	}

	if move.LaidTile != nil {
		move.LaidTile.PlayerName = currentPlayer.Name
		if err := g.LayTile(ctx, currentPlayer.Name, move.LaidTile); err != nil {
			t.Fatalf("could not lay tile: %v", err)
		}
	} else if move.Draw {
		if !g.DrawTile(ctx, currentPlayer.Name) {
			t.Fatal("could not draw tile")
		}
	}

	if r := g.CurrentRound(ctx); r == nil {
		t.Fatalf("Round is done after move (player killed own line): %s", move)
	}

}
