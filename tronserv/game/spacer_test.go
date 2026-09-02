package game

import "testing"

func TestLaySpacerNotStraight(t *testing.T) {
	ctx := t.Context()
	_, r := newLeaderRound(4, "alice")
	err := r.LaySpacer(ctx, &Game{BoardWidth: 21, BoardHeight: 21}, "alice", &Spacer{A: Coord{10, 9}, B: Coord{13, 4}})
	if err != ErrSpacerNotStraight {
		t.Fatalf("diagonal spacer: err=%v, want ErrSpacerNotStraight", err)
	}
}

func TestLaySpacerWrongLength(t *testing.T) {
	ctx := t.Context()
	_, r := newLeaderRound(4, "alice")
	err := r.LaySpacer(ctx, &Game{BoardWidth: 21, BoardHeight: 21}, "alice", &Spacer{A: Coord{10, 9}, B: Coord{10, 7}})
	if err != ErrWrongLengthSpacer {
		t.Fatalf("short spacer: err=%v, want ErrWrongLengthSpacer", err)
	}
}

func TestLaySpacerOccludedPath(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	// Block the path 2 squares up from the spacer's start.
	r.LaidTiles = append(r.LaidTiles, &LaidTile{Tile: Tile{PipsA: 1, PipsB: 1}, Coord: Coord{10, 7}, Orientation: "right"})
	err := r.LaySpacer(ctx, g, "alice", &Spacer{A: Coord{10, 9}, B: Coord{10, 4}})
	if err != ErrTileOccluded {
		t.Fatalf("blocked path: err=%v, want ErrTileOccluded", err)
	}
}

func TestLaySpacerNotStartedOnLine(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	err := r.LaySpacer(ctx, g, "alice", &Spacer{A: Coord{2, 2}, B: Coord{2, 7}})
	if err != ErrSpacerNotStartedOnLine {
		t.Fatalf("spacer not touching any line head: err=%v, want ErrSpacerNotStartedOnLine", err)
	}
}

func TestLaySpacerSuccess(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	spacer := &Spacer{A: Coord{10, 9}, B: Coord{10, 4}}
	if err := r.LaySpacer(ctx, g, "alice", spacer); err != nil {
		t.Fatalf("valid spacer off the leader: %v", err)
	}
	if r.Spacer == nil || *r.Spacer != *spacer {
		t.Errorf("r.Spacer = %v, want %v", r.Spacer, spacer)
	}
}

func TestLaySpacerReversedWhenOnlyFarEndTouchesLine(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	// A doesn't touch the leader's line head; B does. LaySpacer should
	// silently canonicalize by swapping A and B.
	given := &Spacer{A: Coord{10, 4}, B: Coord{10, 9}}
	if err := r.LaySpacer(ctx, g, "alice", given); err != nil {
		t.Fatalf("valid spacer given backwards: %v", err)
	}
	want := Spacer{A: Coord{10, 9}, B: Coord{10, 4}}
	if r.Spacer == nil || *r.Spacer != want {
		t.Errorf("r.Spacer = %v, want %v (reversed)", r.Spacer, want)
	}
}

func TestGameLaySpacerRejectsWhenChickenFooted(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice")
	g.Players[0].ChickenFoot = true
	err := g.LaySpacer(ctx, "alice", &Spacer{A: Coord{10, 9}, B: Coord{10, 4}})
	if err != ErrSpacerNoChickenFoot {
		t.Fatalf("LaySpacer while chicken-footed: err=%v, want ErrSpacerNoChickenFoot", err)
	}
}

func TestGameLaySpacerRejectsWrongTurn(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice", "bob")
	err := g.LaySpacer(ctx, "bob", &Spacer{A: Coord{10, 9}, B: Coord{10, 4}})
	if err != ErrNotYourTurn {
		t.Fatalf("LaySpacer out of turn: err=%v, want ErrNotYourTurn", err)
	}
}

func TestGameLaySpacerClearsWithEmptySpacer(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	r.Spacer = &Spacer{A: Coord{10, 9}, B: Coord{10, 4}}
	if err := g.LaySpacer(ctx, "alice", &Spacer{}); err != nil {
		t.Fatalf("clearing spacer: %v", err)
	}
	if r.Spacer != nil {
		t.Errorf("r.Spacer = %v, want nil after clearing", r.Spacer)
	}
}

// -- CanSetChickenFoot / SetChickenFoot --

func TestCanSetChickenFootSentinel(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	err := r.CanSetChickenFoot(ctx, g, g.Players[0], Coord{-1, -1})
	if err != ErrMustPickChickenFoot {
		t.Fatalf("sentinel coord: err=%v, want ErrMustPickChickenFoot", err)
	}
}

func TestCanSetChickenFootOnLeaderItself(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	leader := r.LaidTiles[0]
	err := r.CanSetChickenFoot(ctx, g, g.Players[0], leader.CoordA())
	if err != ErrMustPickChickenFoot {
		t.Fatalf("coord on the leader tile itself: err=%v, want ErrMustPickChickenFoot", err)
	}
}

func TestCanSetChickenFootNotAdjacent(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	err := r.CanSetChickenFoot(ctx, g, g.Players[0], Coord{2, 2})
	if err != ErrMustPickChickenFoot {
		t.Fatalf("coord far from the leader: err=%v, want ErrMustPickChickenFoot", err)
	}
}

func TestCanSetChickenFootOccluded(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	leader := r.LaidTiles[0]
	blockCoord := leader.CoordA().Up() // (10,9), adjacent to the leader
	r.LaidTiles = append(r.LaidTiles, &LaidTile{Tile: Tile{PipsA: 1, PipsB: 1}, Coord: blockCoord, Orientation: "right"})
	err := r.CanSetChickenFoot(ctx, g, g.Players[0], blockCoord)
	if err != ErrTileOccluded {
		t.Fatalf("occupied coord: err=%v, want ErrTileOccluded", err)
	}
}

func TestCanSetChickenFootDuplicate(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice", "bob")
	shared := Coord{10, 9}
	g.Players[1].ChickenFoot = true
	g.Players[1].ChickenFootCoord = shared
	r := g.Rounds[0]
	err := r.CanSetChickenFoot(ctx, g, g.Players[0], shared)
	if err != ErrBadChickenFoot {
		t.Fatalf("coord already claimed by another footed player: err=%v, want ErrBadChickenFoot", err)
	}
}

func TestCanSetChickenFootSuccess(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	leader := r.LaidTiles[0]
	coord := leader.CoordA().Up()
	if err := r.CanSetChickenFoot(ctx, g, g.Players[0], coord); err != nil {
		t.Fatalf("valid chicken-foot coord: %v", err)
	}
	if err := r.SetChickenFoot(ctx, g, g.Players[0], coord); err != nil {
		t.Fatalf("SetChickenFoot: %v", err)
	}
	if g.Players[0].ChickenFootCoord != coord {
		t.Errorf("ChickenFootCoord = %s, want %s", g.Players[0].ChickenFootCoord, coord)
	}
}

func TestCanSetChickenFootNoOpenAdjacent(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	leader := r.LaidTiles[0]
	coord := leader.CoordA().Up() // (10,9)
	// Occupy all 4 neighbors of the candidate coord so nothing is left
	// open for the player's line to actually continue into. Each blocker
	// is oriented away from coord so its far end doesn't spill back onto
	// coord itself (Neighbors() returns [Up, Down, Left, Right]).
	directions := []string{"up", "down", "left", "right"}
	for i, n := range coord.Neighbors() {
		r.LaidTiles = append(r.LaidTiles, &LaidTile{
			Tile: Tile{PipsA: i, PipsB: i}, Coord: n, Orientation: directions[i],
		})
	}
	err := r.CanSetChickenFoot(ctx, g, g.Players[0], coord)
	if err != ErrNoOpenAdjacent {
		t.Fatalf("boxed-in chicken foot: err=%v, want ErrNoOpenAdjacent", err)
	}
}
