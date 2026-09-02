package game

// Shared builders for synthetic test positions, complementing the
// bug-report fixtures in testdata/*.json. Constructing state directly
// (instead of driving it through Start(), which shuffles a real bag) keeps
// these scenarios small, deterministic, and focused on one rule at a time.

// newLeaderRound builds a Game with the given players (in turn order) and a
// single in-progress round whose only laid tile is a `pips:pips` leader
// double at the board center. Every player's line starts as [leader],
// matching the state right after Game.Start() places the round leader. The
// board is oversized (21x21) so geometry math never runs out of room.
func newLeaderRound(pips int, names ...string) (*Game, *Round) {
	g := &Game{
		Code:        "TESTAA",
		BoardWidth:  21,
		BoardHeight: 21,
		MaxPips:     16,
	}
	for _, n := range names {
		g.Players = append(g.Players, &Player{Name: n})
	}
	leader := &LaidTile{
		Tile:        Tile{PipsA: pips, PipsB: pips},
		Coord:       Coord{X: 10, Y: 10},
		Orientation: "right",
		NextPips:    pips,
	}
	playerLines := map[string][]*LaidTile{}
	for _, n := range names {
		playerLines[n] = []*LaidTile{leader}
	}
	round := &Round{
		LaidTiles:     []*LaidTile{leader},
		PlayerLines:   playerLines,
		HighestLeader: pips,
	}
	g.Rounds = []*Round{round}
	g.Turn = 0
	return g, round
}
