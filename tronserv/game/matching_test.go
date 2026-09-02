package game

import "testing"

// These exercise Round.canPlayOnLine / canPlayOnTileWithoutIndication, the
// low-level adjacency+pips matcher that both LayTile and FindLegalMoves rely
// on. It is unexported but easy to pin down in isolation, which makes it a
// good place to lock in exact behavior before any performance rework of the
// tile-laying path.

func TestCanPlayOnLineNotAdjacent(t *testing.T) {
	r := &Round{}
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	lt := &LaidTile{Tile: Tile{PipsA: 3, PipsB: 4}, Coord: Coord{20, 20}, Orientation: "right"}
	ok, _, err := r.canPlayOnLine(lt, []*LaidTile{last})
	if ok || err != ErrNotAdjacent {
		t.Fatalf("far-away tile: ok=%v err=%v, want ok=false err=ErrNotAdjacent", ok, err)
	}
}

func TestCanPlayOnLineMustMatchPips(t *testing.T) {
	r := &Round{}
	// last occupies (5,5)-(6,5), open end (NextPips) is the 3 at (6,5).
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	// lt is adjacent (its CoordA (7,5) touches last's CoordB (6,5)) but
	// carries neither pip value that would continue the line.
	lt := &LaidTile{Tile: Tile{PipsA: 5, PipsB: 6}, Coord: Coord{7, 5}, Orientation: "right"}
	ok, _, err := r.canPlayOnLine(lt, []*LaidTile{last})
	if ok || err != ErrMustMatchPips {
		t.Fatalf("mismatched pips: ok=%v err=%v, want ok=false err=ErrMustMatchPips", ok, err)
	}
}

func TestCanPlayOnLineSuccess(t *testing.T) {
	r := &Round{}
	// last occupies (5,5)=2 - (6,5)=3, open end is the 3.
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	// lt occupies (7,5)=3 - (8,5)=4: its A touches last's B, and its A pips (3) match NextPips.
	lt := &LaidTile{Tile: Tile{PipsA: 3, PipsB: 4}, Coord: Coord{7, 5}, Orientation: "right"}
	ok, nextPips, err := r.canPlayOnLine(lt, []*LaidTile{last})
	if !ok || err != nil {
		t.Fatalf("valid extension: ok=%v err=%v, want ok=true err=nil", ok, err)
	}
	if nextPips != 4 {
		t.Errorf("nextPips = %d, want 4 (the far end of the newly laid tile)", nextPips)
	}
}

func TestCanPlayOnLineSuccessOtherSide(t *testing.T) {
	r := &Round{}
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	// This time lt's B end (matching pips=3) touches last's B end, via a
	// tile laid "backwards" relative to the previous test: CoordA=(8,5),
	// CoordB=(7,5) i.e. orientation "left", so CoordB (the 3) is what's
	// adjacent to last's open end.
	lt := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 3}, Coord: Coord{8, 5}, Orientation: "left"}
	ok, nextPips, err := r.canPlayOnLine(lt, []*LaidTile{last})
	if !ok || err != nil {
		t.Fatalf("valid extension (B end touches): ok=%v err=%v, want ok=true err=nil", ok, err)
	}
	if nextPips != 4 {
		t.Errorf("nextPips = %d, want 4", nextPips)
	}
}

func TestCanPlayOnLineWrongSide(t *testing.T) {
	r := &Round{}
	// last occupies (5,5)=2 - (6,5)=3, open end (NextPips) is the 3 at (6,5).
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	// lt carries a matching 3 (PipsA), and lt.CoordB happens to be
	// adjacent to last's open end (6,5) -- but the matching pips are on
	// lt.CoordA, which is two squares away and touches nothing. The tile
	// could be played here if the caller flipped it end-for-end first.
	lt := &LaidTile{Tile: Tile{PipsA: 3, PipsB: 7}, Coord: Coord{6, 3}, Orientation: "down"}
	// CoordA=(6,3) pips=3, CoordB=(6,4) pips=7.
	// AA = (5,5).Adj(6,3) = false
	// AB = (5,5).Adj(6,4) = false
	// BA = (6,5).Adj(6,3) = false
	// BB = (6,5).Adj(6,4) = true  <- adjacency exists, but on the unmatched end
	ok, _, err := r.canPlayOnLine(lt, []*LaidTile{last})
	if ok || err != ErrWrongSide {
		t.Fatalf("flipped tile: ok=%v err=%v, want ok=false err=ErrWrongSide", ok, err)
	}
}

func TestCanPlayOnLineDoubleMatch(t *testing.T) {
	r := &Round{}
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	// A double (3:3) extends the line, and its "next" pips are still 3
	// since both ends carry the same value.
	lt := &LaidTile{Tile: Tile{PipsA: 3, PipsB: 3}, Coord: Coord{7, 5}, Orientation: "right"}
	ok, nextPips, err := r.canPlayOnLine(lt, []*LaidTile{last})
	if !ok || err != nil {
		t.Fatalf("double extension: ok=%v err=%v, want ok=true err=nil", ok, err)
	}
	if nextPips != 3 {
		t.Errorf("nextPips = %d, want 3", nextPips)
	}
}

func TestCanPlayOnTileIndicatingMismatch(t *testing.T) {
	r := &Round{}
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	lt := &LaidTile{
		Tile: Tile{PipsA: 3, PipsB: 4}, Coord: Coord{7, 5}, Orientation: "right",
		Indicating: true, Indicated: Tile{PipsA: 9, PipsB: 9}, // claims `last` is a 9:9, which it isn't
	}
	_, _, err := r.canPlayOnTile(lt, last)
	if err != ErrMustMatchPips {
		t.Fatalf("indicated-tile mismatch: err=%v, want ErrMustMatchPips", err)
	}
}

func TestCanPlayOnTileIndicatingMatch(t *testing.T) {
	r := &Round{}
	last := &LaidTile{Tile: Tile{PipsA: 2, PipsB: 3}, Coord: Coord{5, 5}, Orientation: "right", NextPips: 3}
	lt := &LaidTile{
		Tile: Tile{PipsA: 3, PipsB: 4}, Coord: Coord{7, 5}, Orientation: "right",
		Indicating: true, Indicated: Tile{PipsA: 2, PipsB: 3}, // correctly names `last`
	}
	ok, _, err := r.canPlayOnTile(lt, last)
	if !ok || err != nil {
		t.Fatalf("correctly indicated tile: ok=%v err=%v, want ok=true err=nil", ok, err)
	}
}
