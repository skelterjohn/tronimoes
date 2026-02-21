package game

import (
	"encoding/json"
	"testing"
)

var encodedGames = map[string]string{
	"basic_report": `{"created":1771699419,"version":3,"pickup":true,"done":false,"code":"VOWNSL-QHDDRQ","players":[{"name":"skelterjohn","ready":false,"score":0,"hand":[{"pips_a":2,"pips_b":3},{"pips_a":0,"pips_b":6},{"pips_a":0,"pips_b":0},{"pips_a":1,"pips_b":6},{"pips_a":0,"pips_b":4},{"pips_a":1,"pips_b":3},{"pips_a":1,"pips_b":2},{"pips_a":3,"pips_b":6}],"hints":null,"spacer_hints":null,"chicken_foot":false,"dead":false,"just_drew":false,"chicken_foot_coord":{"x":0,"y":0},"chicken_foot_url":"","react_url":"","kills":null}],"turn":0,"rounds":[{"laid_tiles":[{"tile":{"pips_a":3,"pips_b":3},"coord":{"x":2,"y":3},"orientation":"right","player_name":"","next_pips":3,"dead":false,"indicated":null},{"tile":{"pips_a":0,"pips_b":3},"coord":{"x":3,"y":1},"orientation":"down","player_name":"skelterjohn","next_pips":0,"dead":false,"indicated":{"pips_a":-1,"pips_b":-1}}],"spacer":null,"done":false,"history":["skelterjohn laid 3:3","skelterjohn laid 0:3"],"player_lines":{"skelterjohn":[{"tile":{"pips_a":3,"pips_b":3},"coord":{"x":2,"y":3},"orientation":"right","player_name":"","next_pips":3,"dead":false,"indicated":null},{"tile":{"pips_a":0,"pips_b":3},"coord":{"x":3,"y":1},"orientation":"down","player_name":"skelterjohn","next_pips":0,"dead":false,"indicated":{"pips_a":-1,"pips_b":-1}}]},"free_lines":[],"bagless_passes":0,"highest_leader":3}],"bag":[{"pips_a":0,"pips_b":5},{"pips_a":0,"pips_b":1},{"pips_a":0,"pips_b":2},{"pips_a":2,"pips_b":4},{"pips_a":3,"pips_b":5},{"pips_a":5,"pips_b":5},{"pips_a":4,"pips_b":4},{"pips_a":2,"pips_b":6},{"pips_a":4,"pips_b":5},{"pips_a":1,"pips_b":1},{"pips_a":1,"pips_b":5},{"pips_a":2,"pips_b":5},{"pips_a":4,"pips_b":6},{"pips_a":5,"pips_b":6},{"pips_a":3,"pips_b":4},{"pips_a":6,"pips_b":6},{"pips_a":1,"pips_b":4},{"pips_a":2,"pips_b":2}],"board_width":6,"board_height":7,"max_pips":6,"history":["skelterjohn started round 1 - 3:3"]}`,
}

func decodeGame(t *testing.T, label string) *Game {
	encoded := encodedGames[label]
	var game Game
	if err := json.Unmarshal([]byte(encoded), &game); err != nil {
		t.Fatalf("could not unmarshal: %v", err)
	}
	return &game
}

func testLegalMovesContains(t *testing.T, label string, expectedMoves []*LaidTile, expectedSpacers []*Spacer) {
	game := decodeGame(t, label)
	r := game.CurrentRound(t.Context())
	moves, spacers := r.FindLegalMoves(t.Context(), game, game.Players[0])

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
	testLegalMovesContains(t, "basic_report", []*LaidTile{
		{Tile: &Tile{PipsA: 0, PipsB: 4}, Coord: Coord{X: 4, Y: 1}, Orientation: "right"},
		{Tile: &Tile{PipsA: 0, PipsB: 6}, Coord: Coord{X: 3, Y: 0}, Orientation: "right"},
	}, []*Spacer{})
}
