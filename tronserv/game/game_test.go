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
		{Tile: Tile{PipsA: 0, PipsB: 4}, Coord: Coord{X: 4, Y: 1}, Orientation: "right"},
		{Tile: Tile{PipsA: 0, PipsB: 6}, Coord: Coord{X: 3, Y: 0}, Orientation: "right"},
	}, []*Spacer{})
}

func TestChickenFoot(t *testing.T) {
	game := decodeGame(t, "chickenfoot")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: Tile{PipsA: 1, PipsB: 2}, Coord: Coord{X: 3, Y: 8}, Orientation: "up"},
	}, []*Spacer{})
}

func TestAdjacent(t *testing.T) {
	game := decodeGame(t, "adjacent")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: Tile{PipsA: 1, PipsB: 2}, Coord: Coord{X: 2, Y: 7}, Orientation: "right"},
	}, []*Spacer{})
}

func TestRejectedDouble(t *testing.T) {
	game := decodeGame(t, "rejected_double")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: Tile{PipsA: 6, PipsB: 6}, Coord: Coord{X: 3, Y: 3}, Orientation: "down"},
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
		{Tile: Tile{PipsA: 1, PipsB: 6}, Coord: Coord{X: 4, Y: 4}, Orientation: "right"},
	}, []*Spacer{})
}

func TestLeadfoot2(t *testing.T) {
	game := decodeGame(t, "leadfoot2")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: Tile{PipsA: 2, PipsB: 6}, Coord: Coord{X: 5, Y: 3}, Orientation: "down"},
	}, []*Spacer{})
}

func TestPlayfoot(t *testing.T) {
	game := decodeGame(t, "playfoot")
	testLegalMovesContains(t, game, []*LaidTile{
		{Tile: Tile{PipsA: 0, PipsB: 7}, Coord: Coord{X: 1, Y: 2}, Orientation: "right"},
	}, []*Spacer{})

	err := game.LayTile(t.Context(), "TJ.Tice", &LaidTile{
		Tile:        Tile{PipsA: 0, PipsB: 7},
		Coord:       Coord{X: 1, Y: 2},
		Orientation: "right",
		PlayerName:  "TJ.Tice",
		Indicated:   &Tile{PipsA: 2, PipsB: 6},
	})
	if err != nil {
		t.Fatalf("error laying tile: %v", err)
	}
}

func TestDoublesturn(t *testing.T) {
	game := decodeGame(t, "doublesturn")
	if game.Turn != 1 {
		t.Fatalf("doublesturn testdata: expected turn 1, got %d", game.Turn)
	}
	player1Name := game.Players[1].Name
	lt := &LaidTile{
		Tile:        Tile{PipsA: 2, PipsB: 5},
		Coord:       Coord{X: 4, Y: 7},
		Orientation: "down",
		PlayerName:  player1Name,
		Indicated:   &Tile{PipsA: 3, PipsB: 5},
	}
	err := game.LayTile(t.Context(), player1Name, lt)
	if err != nil {
		t.Fatalf("LayTile(2:5 at (4,7) down): %v", err)
	}
	if game.Turn != 0 {
		t.Errorf("after player 1 lays %s, expected turn 0, got %d", lt, game.Turn)
	}
}

func TestSpacerPlacementRejected(t *testing.T) {
	game := decodeGame(t, "spacerplacement")
	name := game.Players[game.Turn].Name
	// 6:6 at (0,4) "right" would start a free line but is under the spacer (spacer runs (0,4)-(5,4)); must return ErrFreeFromSpacer.
	err := game.LayTile(t.Context(), name, &LaidTile{
		Tile:        Tile{PipsA: 6, PipsB: 6},
		Coord:       Coord{X: 0, Y: 4},
		Orientation: "right",
		PlayerName:  name,
	})
	if err != ErrFreeFromSpacer {
		t.Fatalf("LayTile under spacer: got err %v, want ErrFreeFromSpacer", err)
	}
}
