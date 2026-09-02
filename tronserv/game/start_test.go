package game

import "testing"

// Start() shuffles a fresh bag with the global math/rand source, which makes
// most of it non-deterministic to test directly. The game code carves out
// two special dev/test codes, "AAAAAA" and "BBBBBB", that install a fixed,
// unshuffled bag instead -- these give us full determinism for exercising
// the round-leader search and initial deal.

func TestStartNotEnoughPlayers(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ZZZZZZ")
	if err := g.Start(ctx, "nobody"); err != ErrGameNotEnoughPlayers {
		t.Fatalf("Start with 0 players = %v, want ErrGameNotEnoughPlayers", err)
	}
}

func TestStartPreviousRoundNotDone(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ZZZZZZ")
	g.AddPlayer(ctx, &Player{Name: "alice"})
	g.Rounds = append(g.Rounds, &Round{Done: false})
	if err := g.Start(ctx, "alice"); err != ErrGamePreviousRoundNotDone {
		t.Fatalf("Start with unfinished round = %v, want ErrGamePreviousRoundNotDone", err)
	}
}

func TestStartGameOver(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ZZZZZZ")
	g.AddPlayer(ctx, &Player{Name: "alice"})
	// The game ends once a round led by the 0:0 double has completed.
	g.Rounds = append(g.Rounds, &Round{
		Done:      true,
		LaidTiles: []*LaidTile{{Tile: Tile{PipsA: 0, PipsB: 0}}},
	})
	if err := g.Start(ctx, "alice"); err != ErrGameOver {
		t.Fatalf("Start after 0:0 round = %v, want ErrGameOver", err)
	}
}

func TestStartReadinessGating(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "AAAAAA")
	g.AddPlayer(ctx, &Player{Name: "alice"})
	g.AddPlayer(ctx, &Player{Name: "bob"})

	if err := g.Start(ctx, "alice"); err != nil {
		t.Fatalf("Start(alice): %v", err)
	}
	if !g.Players[0].Ready {
		t.Error("alice should be marked Ready")
	}
	if g.Players[1].Ready {
		t.Error("bob should not be Ready yet")
	}
	if len(g.Rounds) != 0 {
		t.Fatalf("round should not start until everyone is ready, got %d rounds", len(g.Rounds))
	}

	if err := g.Start(ctx, "bob"); err != nil {
		t.Fatalf("Start(bob): %v", err)
	}
	if len(g.Rounds) != 1 {
		t.Fatalf("expected round to start once everyone is ready, got %d rounds", len(g.Rounds))
	}
	// Ready flags reset for the next round once this one begins.
	if g.Players[0].Ready || g.Players[1].Ready {
		t.Error("Ready flags should be reset once the round starts")
	}
}

// TestStartDealAndLeaderAAAAAA pins down the exact deal and round-leader
// selection for the fixed "AAAAAA" bag with 2 players (MaxPips=7), so the
// whole Start() pipeline -- bag construction, dealing, leader search,
// leader tile placement, turn advancement -- has a byte-exact regression
// test.
func TestStartDealAndLeaderAAAAAA(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "AAAAAA")
	g.AddPlayer(ctx, &Player{Name: "alice"})
	g.AddPlayer(ctx, &Player{Name: "bob"})

	if err := g.Start(ctx, "alice"); err != nil {
		t.Fatalf("Start(alice): %v", err)
	}
	if err := g.Start(ctx, "bob"); err != nil {
		t.Fatalf("Start(bob): %v", err)
	}

	if len(g.Rounds) != 1 {
		t.Fatalf("expected 1 round, got %d", len(g.Rounds))
	}
	round := g.Rounds[0]
	if len(round.LaidTiles) != 1 {
		t.Fatalf("expected 1 laid tile (the leader), got %d", len(round.LaidTiles))
	}
	leader := round.LaidTiles[0]
	// bob's initial 10-tile hand (from the fixed AAAAAA deck at
	// MaxPips=7) is the first double found scanning down from MaxPips:
	// 7:7 and 6:6 are still in the bag, but 5:5 is in bob's hand.
	if leader.Tile != (Tile{PipsA: 5, PipsB: 5}) {
		t.Fatalf("leader tile = %v, want 5:5", leader.Tile)
	}
	wantX, wantY := g.BoardWidth/2-1, g.BoardHeight/2
	if leader.Coord != (Coord{wantX, wantY}) {
		t.Errorf("leader coord = %s, want (%d,%d)", leader.Coord, wantX, wantY)
	}

	// bob played the leader tile, so it's gone from his hand; turn then
	// advances past bob back to alice.
	bob := g.GetPlayer(ctx, "bob")
	for _, tile := range bob.Hand {
		if tile == (Tile{PipsA: 5, PipsB: 5}) {
			t.Error("bob's hand should no longer contain the 5:5 leader tile")
		}
	}
	if len(bob.Hand) != 9 {
		t.Errorf("bob's hand size = %d, want 9 (10 dealt - 1 played)", len(bob.Hand))
	}
	alice := g.GetPlayer(ctx, "alice")
	if len(alice.Hand) != 10 {
		t.Errorf("alice's hand size = %d, want 10 (untouched)", len(alice.Hand))
	}
	if g.Players[g.Turn].Name != "alice" {
		t.Errorf("turn = %s, want alice", g.Players[g.Turn].Name)
	}

	// The leader was found in the initial deal, so no extra bag draws
	// happened: 22 fixed tiles - 20 dealt (10 each) = 2 left.
	if len(g.Bag) != 2 {
		t.Errorf("bag size = %d, want 2", len(g.Bag))
	}
}

func TestStartDealAndLeaderBBBBBB(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "BBBBBB")
	g.AddPlayer(ctx, &Player{Name: "alice"})

	if err := g.Start(ctx, "alice"); err != nil {
		t.Fatalf("Start(alice): %v", err)
	}
	round := g.Rounds[0]
	leader := round.LaidTiles[0]
	// The fixed BBBBBB deck at MaxPips=6 gives alice 3:3 as the highest
	// double in her initial 10-tile hand (6:6, 5:5, 4:4 are all still in
	// the bag).
	if leader.Tile != (Tile{PipsA: 3, PipsB: 3}) {
		t.Fatalf("leader tile = %v, want 3:3", leader.Tile)
	}
}
