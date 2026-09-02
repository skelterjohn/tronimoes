package game

import "testing"

func TestPassRequiresDrawWhenBagNonEmpty(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice", "bob")
	g.Bag = []Tile{{PipsA: 1, PipsB: 1}}
	if err := g.Pass(ctx, "alice", -1, -1); err != ErrMustDrawTile {
		t.Fatalf("Pass without drawing (bag non-empty) = %v, want ErrMustDrawTile", err)
	}
}

func TestPassAfterDrawing(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	g.Bag = []Tile{{PipsA: 1, PipsB: 1}}
	g.Players[0].JustDrew = true

	// alice's line is still just the leader (length 1), so passing
	// requires a valid chicken-foot coordinate.
	if err := g.Pass(ctx, "alice", 10, 9); err != nil {
		t.Fatalf("Pass after drawing: %v", err)
	}
	if g.Players[0].JustDrew {
		t.Error("JustDrew should be cleared after passing")
	}
	if g.Turn != 1 {
		t.Errorf("Turn = %d, want 1 (advanced to bob)", g.Turn)
	}
	if r.BaglessPasses != 0 {
		t.Errorf("BaglessPasses = %d, want 0 (bag was non-empty)", r.BaglessPasses)
	}
	if !g.Players[0].ChickenFoot {
		t.Error("alice should be chicken-footed after her first pass")
	}
}

func TestPassWithEmptyBagNeedsNoDraw(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	g.Bag = nil

	if err := g.Pass(ctx, "alice", 10, 9); err != nil {
		t.Fatalf("Pass with empty bag: %v", err)
	}
	if r.BaglessPasses != 1 {
		t.Errorf("BaglessPasses = %d, want 1", r.BaglessPasses)
	}
}

func TestPassUnknownPlayer(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice")
	if err := g.Pass(ctx, "mallory", 10, 9); err != ErrNoSuchPlayer {
		t.Fatalf("Pass(unknown player) = %v, want ErrNoSuchPlayer", err)
	}
}

func TestPassChickenFootFromExistingLine(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	g.Bag = nil
	leader := r.LaidTiles[0]

	// Give alice a second tile in her line, extending off the leader's
	// open end (CoordB, since NextPips==PipsA for the leader double).
	ext := &LaidTile{
		Tile: Tile{PipsA: 4, PipsB: 7}, Coord: leader.CoordB(), Orientation: "right",
		NextPips: 7, PlayerName: "alice",
	}
	r.PlayerLines["alice"] = append(r.PlayerLines["alice"], ext)
	r.LaidTiles = append(r.LaidTiles, ext)

	if err := g.Pass(ctx, "alice", -1, -1); err != nil {
		t.Fatalf("Pass with an established line: %v", err)
	}
	if !g.Players[0].ChickenFoot {
		t.Fatal("alice should be chicken-footed after passing")
	}
	// NextPips (7) equals ext.Tile.PipsB, so per the code the foot lands
	// on ext.CoordA() (the near end), not CoordB.
	if g.Players[0].ChickenFootCoord != ext.CoordA() {
		t.Errorf("ChickenFootCoord = %s, want %s (ext.CoordA)", g.Players[0].ChickenFootCoord, ext.CoordA())
	}
}

func TestPassInvalidChickenFootCoord(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice", "bob")
	g.Bag = nil
	// Sentinel "not chosen" coordinate.
	if err := g.Pass(ctx, "alice", -1, -1); err != ErrMustPickChickenFoot {
		t.Fatalf("Pass with no chicken-foot coord chosen = %v, want ErrMustPickChickenFoot", err)
	}
}

func TestPassStalemate(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	g.Bag = nil

	if err := g.Pass(ctx, "alice", 10, 9); err != nil {
		t.Fatalf("alice's pass: %v", err)
	}
	if r.Done {
		t.Fatal("round should not be done after only 1 of 2 players passed")
	}
	if err := g.Pass(ctx, "bob", 11, 9); err != nil {
		t.Fatalf("bob's pass: %v", err)
	}
	if !r.Done {
		t.Fatal("round should be Done once every player has passed on an empty bag (stalemate)")
	}
	for _, lt := range r.LaidTiles {
		if !lt.Dead {
			t.Errorf("laid tile %s should be marked Dead after a stalemate", lt)
		}
	}
	for _, p := range g.Players {
		if !p.Dead {
			t.Errorf("player %s should be marked Dead after a stalemate", p.Name)
		}
		if p.ChickenFoot {
			t.Errorf("player %s should have ChickenFoot cleared after a stalemate", p.Name)
		}
	}
}

func TestDrawTile(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice")
	g.Bag = []Tile{{PipsA: 2, PipsB: 5}, {PipsA: 1, PipsB: 1}}
	alice := g.Players[0]
	startHandLen := len(alice.Hand)

	if ok := g.DrawTile(ctx, "alice"); !ok {
		t.Fatal("DrawTile should succeed")
	}
	if len(alice.Hand) != startHandLen+1 {
		t.Errorf("hand size = %d, want %d", len(alice.Hand), startHandLen+1)
	}
	if alice.Hand[len(alice.Hand)-1] != (Tile{PipsA: 2, PipsB: 5}) {
		t.Errorf("drew %v, want the top of the bag (2:5)", alice.Hand[len(alice.Hand)-1])
	}
	if len(g.Bag) != 1 {
		t.Errorf("bag size = %d, want 1", len(g.Bag))
	}
	if !alice.JustDrew {
		t.Error("JustDrew should be true after drawing")
	}

	// Drawing again before playing/passing is rejected outright.
	if ok := g.DrawTile(ctx, "alice"); ok {
		t.Fatal("DrawTile should return false when the player has already just drawn")
	}
	if len(alice.Hand) != startHandLen+1 {
		t.Errorf("second draw should not add a tile: hand size = %d, want %d", len(alice.Hand), startHandLen+1)
	}

	if ok := g.DrawTile(ctx, "nobody"); ok {
		t.Error("DrawTile(unknown player) should return false")
	}
}

// TestPassWithNoCurrentRoundPanics pins down a real bug: Pass() fetches the
// current round and unconditionally increments/resets round.BaglessPasses
// before it ever checks the round for nil (the nil check exists, but only
// guards the later round.Spacer assignment). Calling Pass when there is no
// active round -- e.g. between rounds, or after the game is done -- crashes
// instead of returning ErrRoundNotStarted. This test documents today's
// actual behavior so it isn't silently changed by unrelated refactoring;
// flagged separately as a bug worth fixing (return ErrRoundNotStarted as
// soon as CurrentRound is nil, before touching round.BaglessPasses).
func TestPassWithNoCurrentRoundPanics(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ZZZZZZ")
	g.AddPlayer(ctx, &Player{Name: "alice"})
	// No rounds at all, so CurrentRound(ctx) is nil.
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Pass() with no current round to panic (known bug); if this now passes cleanly, replace this test with one asserting ErrRoundNotStarted")
		}
	}()
	g.Pass(ctx, "alice", 0, 0)
}

func TestDrawTileEmptyBag(t *testing.T) {
	ctx := t.Context()
	g, _ := newLeaderRound(4, "alice")
	g.Bag = nil
	alice := g.Players[0]
	startHandLen := len(alice.Hand)

	if ok := g.DrawTile(ctx, "alice"); !ok {
		t.Fatal("DrawTile with an empty bag still reports success (nothing to draw)")
	}
	if len(alice.Hand) != startHandLen {
		t.Errorf("hand size changed with an empty bag: got %d, want %d", len(alice.Hand), startHandLen)
	}
	if alice.JustDrew {
		t.Error("JustDrew should stay false when the bag is empty")
	}
}
