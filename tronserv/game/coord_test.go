package game

import "testing"

func TestCoordAdj(t *testing.T) {
	cases := []struct {
		a, b Coord
		want bool
	}{
		{Coord{0, 0}, Coord{0, 1}, true},
		{Coord{0, 0}, Coord{0, -1}, true},
		{Coord{0, 0}, Coord{1, 0}, true},
		{Coord{0, 0}, Coord{-1, 0}, true},
		{Coord{0, 0}, Coord{0, 0}, false},  // not adjacent to itself
		{Coord{0, 0}, Coord{1, 1}, false},  // diagonal is not adjacent
		{Coord{0, 0}, Coord{0, 2}, false},  // two away is not adjacent
		{Coord{5, 5}, Coord{5, 6}, true},
	}
	for _, c := range cases {
		if got := c.a.Adj(c.b); got != c.want {
			t.Errorf("%s.Adj(%s) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestCoordOrientationTo(t *testing.T) {
	cases := []struct {
		a, b Coord
		want string
	}{
		{Coord{0, 0}, Coord{0, 1}, "down"},
		{Coord{0, 1}, Coord{0, 0}, "up"},
		{Coord{0, 0}, Coord{1, 0}, "right"},
		{Coord{1, 0}, Coord{0, 0}, "left"},
		{Coord{0, 0}, Coord{1, 1}, ""}, // neither same row nor column
		{Coord{0, 0}, Coord{0, 0}, "up"}, // same point: X==o.X and !(Y<o.Y) so "up"
	}
	for _, c := range cases {
		if got := c.a.OrientationTo(c.b); got != c.want {
			t.Errorf("%s.OrientationTo(%s) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestCoordOrientation(t *testing.T) {
	c := Coord{5, 5}
	cases := []struct {
		o    string
		want Coord
	}{
		{"up", Coord{5, 4}},
		{"down", Coord{5, 6}},
		{"left", Coord{4, 5}},
		{"right", Coord{6, 5}},
		{"sideways", Coord{5, 5}}, // unrecognized orientation is a no-op
	}
	for _, tc := range cases {
		if got := c.Orientation(tc.o); got != tc.want {
			t.Errorf("Coord{5,5}.Orientation(%q) = %s, want %s", tc.o, got, tc.want)
		}
	}
}

func TestCoordNeighbors(t *testing.T) {
	got := Coord{5, 5}.Neighbors()
	want := []Coord{{5, 4}, {5, 6}, {4, 5}, {6, 5}}
	if len(got) != len(want) {
		t.Fatalf("Neighbors() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Neighbors()[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestLaidTileCoordAB(t *testing.T) {
	lt := LaidTile{Tile: Tile{PipsA: 3, PipsB: 4}, Coord: Coord{2, 3}, Orientation: "right"}
	if lt.CoordA() != (Coord{2, 3}) {
		t.Errorf("CoordA() = %s, want (2,3)", lt.CoordA())
	}
	if lt.CoordB() != (Coord{3, 3}) {
		t.Errorf("CoordB() = %s, want (3,3)", lt.CoordB())
	}
}

func TestLaidTileReverse(t *testing.T) {
	lt := LaidTile{Tile: Tile{PipsA: 3, PipsB: 4}, Coord: Coord{2, 3}, Orientation: "right", NextPips: 4}
	rev := lt.Reverse()
	if rev.Coord != (Coord{3, 3}) {
		t.Errorf("Reverse().Coord = %s, want (3,3)", rev.Coord)
	}
	if rev.Orientation != "left" {
		t.Errorf("Reverse().Orientation = %q, want %q", rev.Orientation, "left")
	}
	// Reversing should preserve the physical footprint of the tile: the
	// reversed tile's A/B coords are the original's B/A coords.
	if rev.CoordA() != lt.CoordB() {
		t.Errorf("rev.CoordA() = %s, want %s (= original CoordB)", rev.CoordA(), lt.CoordB())
	}
	if rev.CoordB() != lt.CoordA() {
		t.Errorf("rev.CoordB() = %s, want %s (= original CoordA)", rev.CoordB(), lt.CoordA())
	}

	orientations := map[string]string{"up": "down", "down": "up", "left": "right", "right": "left"}
	for o, want := range orientations {
		lt := LaidTile{Orientation: o}
		if got := lt.Reverse().Orientation; got != want {
			t.Errorf("Reverse() of orientation %q = %q, want %q", o, got, want)
		}
	}
}

func TestPlayerHasRoundLeader(t *testing.T) {
	p := &Player{Hand: []Tile{{PipsA: 1, PipsB: 2}, {PipsA: 4, PipsB: 4}, {PipsA: 0, PipsB: 3}}}
	if !p.HasRoundLeader(4) {
		t.Errorf("HasRoundLeader(4) = false, want true (hand has 4:4)")
	}
	if p.HasRoundLeader(1) {
		t.Errorf("HasRoundLeader(1) = true, want false (1:2 is not a double)")
	}
	if p.HasRoundLeader(9) {
		t.Errorf("HasRoundLeader(9) = true, want false (no such tile in hand)")
	}
}
