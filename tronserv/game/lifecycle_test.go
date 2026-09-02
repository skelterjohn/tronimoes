package game

import "testing"

func TestNewGame(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ABCDEF")
	if g.Code != "ABCDEF" {
		t.Errorf("Code = %q, want %q", g.Code, "ABCDEF")
	}
	if g.Version != 0 {
		t.Errorf("Version = %d, want 0", g.Version)
	}
	if g.Done {
		t.Error("new game should not be Done")
	}
	if len(g.Players) != 0 {
		t.Errorf("new game should have no players, got %d", len(g.Players))
	}
}

func TestAdjustBoardAndPips(t *testing.T) {
	ctx := t.Context()
	cases := []struct {
		players             int
		width, height, pips int
	}{
		{1, 6, 7, 6},
		{2, 8, 9, 7},
		{3, 10, 11, 8},
		{4, 12, 13, 10},
		{5, 14, 15, 11},
		{6, 16, 17, 12},
	}
	for _, c := range cases {
		g := &Game{}
		for i := 0; i < c.players; i++ {
			g.Players = append(g.Players, &Player{Name: string(rune('A' + i))})
		}
		g.AdjustBoardAndPips(ctx)
		if g.BoardWidth != c.width || g.BoardHeight != c.height || g.MaxPips != c.pips {
			t.Errorf("%d players: got (w=%d,h=%d,pips=%d), want (w=%d,h=%d,pips=%d)",
				c.players, g.BoardWidth, g.BoardHeight, g.MaxPips, c.width, c.height, c.pips)
		}
	}
}

func TestAddPlayer(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ABCDEF")

	if err := g.AddPlayer(ctx, &Player{Name: "alice", Score: 99}); err != nil {
		t.Fatalf("AddPlayer(alice): %v", err)
	}
	if len(g.Players) != 1 {
		t.Fatalf("expected 1 player, got %d", len(g.Players))
	}
	if g.Players[0].Score != 0 {
		t.Errorf("AddPlayer should reset Score to 0, got %d", g.Players[0].Score)
	}
	// Board should have been resized for 1 player.
	if g.BoardWidth != 6 {
		t.Errorf("BoardWidth after 1 player = %d, want 6", g.BoardWidth)
	}

	if err := g.AddPlayer(ctx, &Player{Name: "alice"}); err != ErrPlayerAlreadyInGame {
		t.Errorf("AddPlayer(duplicate alice) = %v, want ErrPlayerAlreadyInGame", err)
	}

	for _, name := range []string{"bob", "carol", "dave", "eve"} {
		if err := g.AddPlayer(ctx, &Player{Name: name}); err != nil {
			t.Fatalf("AddPlayer(%s): %v", name, err)
		}
	}
	if len(g.Players) != 5 {
		t.Fatalf("expected 5 players, got %d", len(g.Players))
	}
	if err := g.AddPlayer(ctx, &Player{Name: "frank"}); err != nil {
		t.Fatalf("AddPlayer(frank, 6th player): %v", err)
	}
	if err := g.AddPlayer(ctx, &Player{Name: "grace"}); err != ErrGameTooManyPlayers {
		t.Errorf("AddPlayer(7th player) = %v, want ErrGameTooManyPlayers", err)
	}
}

func TestAddPlayerAfterGameStarted(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ABCDEF")
	if err := g.AddPlayer(ctx, &Player{Name: "alice"}); err != nil {
		t.Fatalf("AddPlayer(alice): %v", err)
	}
	g.Rounds = append(g.Rounds, &Round{})
	if err := g.AddPlayer(ctx, &Player{Name: "bob"}); err != ErrGameAlreadyStarted {
		t.Errorf("AddPlayer after round started = %v, want ErrGameAlreadyStarted", err)
	}
}

func TestLeaveOrQuitBeforeRoundsStarted(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ABCDEF")
	for _, name := range []string{"alice", "bob", "carol"} {
		if err := g.AddPlayer(ctx, &Player{Name: name}); err != nil {
			t.Fatalf("AddPlayer(%s): %v", name, err)
		}
	}
	// Board should be sized for 3 players before anyone leaves.
	if g.BoardWidth != 10 {
		t.Fatalf("BoardWidth for 3 players = %d, want 10", g.BoardWidth)
	}

	if quit := g.LeaveOrQuit(ctx, "nobody"); quit {
		t.Error("LeaveOrQuit(nobody) = true, want false")
	}
	if len(g.Players) != 3 {
		t.Errorf("player count changed after quitting a nonexistent player: %d", len(g.Players))
	}

	if quit := g.LeaveOrQuit(ctx, "bob"); !quit {
		t.Error("LeaveOrQuit(bob) = false, want true")
	}
	if len(g.Players) != 2 {
		t.Fatalf("expected 2 players after bob leaves, got %d", len(g.Players))
	}
	for _, p := range g.Players {
		if p.Name == "bob" {
			t.Error("bob is still in the player list")
		}
	}
	// Board should have been resized down to 2 players.
	if g.BoardWidth != 8 {
		t.Errorf("BoardWidth after leaving down to 2 players = %d, want 8", g.BoardWidth)
	}
	if g.Done {
		t.Error("game should not be Done while players remain")
	}

	g.LeaveOrQuit(ctx, "alice")
	if g.Done {
		t.Error("game should not be Done with 1 player remaining")
	}
	g.LeaveOrQuit(ctx, "carol")
	if !g.Done {
		t.Error("game should be Done once all players have left")
	}
}

func TestLeaveOrQuitAfterRoundsStarted(t *testing.T) {
	ctx := t.Context()
	g := NewGame(ctx, "ABCDEF")
	for _, name := range []string{"alice", "bob"} {
		g.AddPlayer(ctx, &Player{Name: name})
	}
	g.Rounds = append(g.Rounds, &Round{})

	if quit := g.LeaveOrQuit(ctx, "nobody"); quit {
		t.Error("LeaveOrQuit(nobody) after round started = true, want false")
	}
	if g.Done {
		t.Error("game should not be Done after a no-op quit")
	}

	if quit := g.LeaveOrQuit(ctx, "alice"); !quit {
		t.Error("LeaveOrQuit(alice) after round started = false, want true")
	}
	if !g.Done {
		t.Error("game should be Done once a player quits mid-round")
	}
	// Quitting mid-round doesn't remove the player from the roster (scores
	// still need to be attributed), it just ends the game.
	if len(g.Players) != 2 {
		t.Errorf("quitting mid-round should not remove the player, got %d players", len(g.Players))
	}
}

func TestLastRoundLeader(t *testing.T) {
	ctx := t.Context()
	g := &Game{MaxPips: 9}
	if got := g.LastRoundLeader(ctx); got != 10 {
		t.Errorf("LastRoundLeader with no rounds = %d, want MaxPips+1 = 10", got)
	}

	g.Rounds = append(g.Rounds, &Round{
		LaidTiles: []*LaidTile{{Tile: Tile{PipsA: 5, PipsB: 5}}},
	})
	if got := g.LastRoundLeader(ctx); got != 5 {
		t.Errorf("LastRoundLeader = %d, want 5", got)
	}
}
