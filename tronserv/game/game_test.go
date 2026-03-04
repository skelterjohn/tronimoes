package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func decodeGame(t *testing.T, label string) *Game {
	path := filepath.Join("testdata", label+".json")
	encoded, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", path, err)
	}
	var game Game
	if err := json.Unmarshal(encoded, &game); err != nil {
		t.Fatalf("could not unmarshal: %v", err)
	}
	return &game
}

func testLegalMovesContains(t *testing.T, game *Game, expectedMoves []*LaidTile, expectedSpacers []*Spacer) {
	r := game.CurrentRound(t.Context())
	moves, spacers := r.FindLegalMoves(t.Context(), game, game.Players[game.Turn])

	foundMoveStrings := make(map[string]bool)
	for _, move := range moves {
		foundMoveStrings[move.String()] = true
		t.Logf("found move: %s", move.String())
	}
	for _, move := range expectedMoves {
		if !foundMoveStrings[move.String()] {
			t.Errorf("move %s not found", move.String())
		}
	}
	foundSpacerStrings := make(map[string]bool)
	for _, spacer := range spacers {
		foundSpacerStrings[spacer.String()] = true
		t.Logf("found spacer: %s", spacer.String())
	}
	for _, spacer := range expectedSpacers {
		if !foundSpacerStrings[spacer.String()] {
			t.Errorf("spacer %s not found", spacer.String())
		}
	}
}

func TestBasicReport(t *testing.T) {
	game := decodeGame(t, "basic_report")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 0, PipsB: 4}, Coord: Coord{X: 4, Y: 1}, Orientation: "right"},
		{Tile: &Tile{PipsA: 0, PipsB: 6}, Coord: Coord{X: 3, Y: 0}, Orientation: "right"},
	}, []*Spacer{})
}

func TestChickenFoot(t *testing.T) {
	game := decodeGame(t, "chickenfoot")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 1, PipsB: 2}, Coord: Coord{X: 3, Y: 8}, Orientation: "up"},
	}, []*Spacer{})
}

func TestAdjacent(t *testing.T) {
	game := decodeGame(t, "adjacent")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 1, PipsB: 2}, Coord: Coord{X: 2, Y: 7}, Orientation: "right"},
	}, []*Spacer{})
}

func TestRejectedDouble(t *testing.T) {
	game := decodeGame(t, "rejected_double")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 6, PipsB: 6}, Coord: Coord{X: 3, Y: 3}, Orientation: "down"},
	}, []*Spacer{})
}

func TestFreeLine(t *testing.T) {
	game := decodeGame(t, "freeline")
	testLegalMovesContains(t, game, []*LaidTile{}, []*Spacer{{
		A: Coord{X: 1, Y: 0},
		B: Coord{X: 6, Y: 0},
	}})
}

func TestLeadfoot(t *testing.T) {
	game := decodeGame(t, "leadfoot")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 1, PipsB: 6}, Coord: Coord{X: 4, Y: 4}, Orientation: "right"},
	}, []*Spacer{})
}

func TestLeadfoot2(t *testing.T) {
	game := decodeGame(t, "leadfoot2")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 2, PipsB: 6}, Coord: Coord{X: 5, Y: 3}, Orientation: "down"},
	}, []*Spacer{})
}

func TestPlayfoot(t *testing.T) {
	game := decodeGame(t, "playfoot")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: &Tile{PipsA: 0, PipsB: 7}, Coord: Coord{X: 1, Y: 2}, Orientation: "right"},
	}, []*Spacer{})
}
