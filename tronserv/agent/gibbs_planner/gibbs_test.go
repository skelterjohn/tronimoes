package gibbs_planner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
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
		EvalDecay:          0.9,
	}

	// Previous game has no rounds so Update runs createInitialGuesses.
	previousGame := &game.Game{
		Players: g.Players,
		MaxPips: g.MaxPips,
	}

	gp.Update(ctx, previousGame, g)

	move := gp.GetMove(ctx, g, currentPlayer)
	expectedMove := types.Move{
		LaidTile: &game.LaidTile{
			Tile: &game.Tile{
				PipsA: 3,
				PipsB: 4,
			},
			Coord:       game.Coord{X: 3, Y: 2},
			Orientation: "down",
			Indicated:   nil,
		},
	}
	if move.String() != expectedMove.String() {
		t.Errorf("wrong move: got %s != want %s", move, expectedMove)
	}
}
