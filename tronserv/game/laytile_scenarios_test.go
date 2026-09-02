package game

import "testing"

func TestFindHints(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	alice := g.Players[0]
	// A tile that can legally extend the leader, and one that can't (no 4
	// or 9 pip to match the leader's open end).
	alice.Hand = []Tile{{PipsA: 4, PipsB: 6}, {PipsA: 1, PipsB: 2}}

	r.FindHints(ctx, g, alice)

	if len(alice.Hints) != len(alice.Hand) {
		t.Fatalf("len(Hints) = %d, want %d (one entry per hand tile)", len(alice.Hints), len(alice.Hand))
	}
	if len(alice.Hints[0]) == 0 {
		t.Error("the 4:6 tile should have at least one legal placement hinted")
	}
	if len(alice.Hints[1]) != 0 {
		t.Errorf("the 1:2 tile has no legal placement, want no hints, got %v", alice.Hints[1])
	}
}

// -- Round-ending win conditions --

func TestLayTileWinByEmptyingHand(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	leader := r.LaidTiles[0]

	// bob has already extended his own line once, and still holds tiles
	// (an empty bob.Hand would spuriously also satisfy the "emptied their
	// hand" win check below, since both players are still "living").
	bobExt := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 7}, Coord: Coord{12, 10}, Orientation: "right", NextPips: 7, PlayerName: "bob"}
	r.PlayerLines["bob"] = append(r.PlayerLines["bob"], bobExt)
	r.LaidTiles = append(r.LaidTiles, bobExt)
	g.Players[1].Hand = []Tile{{PipsA: 1, PipsB: 2}}

	alice := g.Players[0]
	alice.Hand = []Tile{{PipsA: 4, PipsB: 5}} // her only tile

	err := g.LayTile(ctx, "alice", &LaidTile{Tile: Tile{PipsA: 4, PipsB: 5}, Coord: Coord{10, 9}, Orientation: "up"})
	if err != nil {
		t.Fatalf("LayTile (last tile in hand): %v", err)
	}
	if !r.Done {
		t.Fatal("round should be Done once a player empties their hand")
	}
	if alice.Score != 2 {
		t.Errorf("alice.Score = %d, want 2 (efficiency win)", alice.Score)
	}
	bob := g.Players[1]
	if !bob.Dead {
		t.Error("bob should be marked Dead when someone else wins by emptying their hand")
	}
	if !leader.Dead {
		t.Error("the leader tile (PlayerName=\"\") should be killed too, only the winner's own tiles survive")
	}
	if !bobExt.Dead {
		t.Error("bob's tile should be killed since he isn't the winner")
	}
	if g.Done {
		t.Error("the game itself should not be over (leader was 4:4, not 0:0)")
	}
}

func TestLayTileGameEndsAfterZeroLeaderRound(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(0, "alice")
	alice := g.Players[0]
	alice.Hand = []Tile{{PipsA: 0, PipsB: 3}}

	err := g.LayTile(ctx, "alice", &LaidTile{Tile: Tile{PipsA: 0, PipsB: 3}, Coord: Coord{10, 9}, Orientation: "up"})
	if err != nil {
		t.Fatalf("LayTile: %v", err)
	}
	if !r.Done {
		t.Fatal("round should be done")
	}
	if !g.Done {
		t.Error("the whole game should end once a round led by the 0:0 double finishes")
	}
}

func TestLayTileAttritionWin(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob", "carol")
	leader := r.LaidTiles[0]

	// bob and carol already died earlier in this round (each had played at
	// least one tile before dying, so their lines are already length>1 --
	// this keeps Game.LayTile's blocking-feet fairness check out of the
	// picture, since it only protects players who haven't started yet).
	bobExt := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 8}, Coord: Coord{12, 10}, Orientation: "right", NextPips: 8, PlayerName: "bob"}
	r.PlayerLines["bob"] = append(r.PlayerLines["bob"], bobExt)
	r.LaidTiles = append(r.LaidTiles, bobExt)
	g.Players[1].Dead = true

	carolExt := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 9}, Coord: Coord{10, 11}, Orientation: "down", NextPips: 9, PlayerName: "carol"}
	r.PlayerLines["carol"] = append(r.PlayerLines["carol"], carolExt)
	r.LaidTiles = append(r.LaidTiles, carolExt)
	g.Players[2].Dead = true

	alice := g.Players[0]
	alice.Hand = []Tile{{PipsA: 4, PipsB: 5}, {PipsA: 9, PipsB: 9}} // 2 tiles, so this play doesn't also win by emptying her hand

	err := g.LayTile(ctx, "alice", &LaidTile{Tile: Tile{PipsA: 4, PipsB: 5}, Coord: Coord{10, 9}, Orientation: "up"})
	if err != nil {
		t.Fatalf("LayTile: %v", err)
	}
	if !r.Done {
		t.Fatal("round should be Done once only one living player remains")
	}
	if alice.Score != 2 {
		t.Errorf("alice.Score = %d, want 2 (attrition win)", alice.Score)
	}
	if leader.Dead {
		t.Error("attrition wins don't retroactively kill tiles the way an efficiency win does")
	}
}

// -- findOuroboros: consuming a footed opponent's exposed line head --

func TestFindOuroborosSelfConsumingOwnLine(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	alice, bob := g.Players[0], g.Players[1]

	// bob extended his line once and is now chicken-footed, waiting for a
	// 6 to land next to (10,8).
	bobTile := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 6}, Coord: Coord{10, 9}, Orientation: "up", NextPips: 6, PlayerName: "bob"}
	r.PlayerLines["bob"] = append(r.PlayerLines["bob"], bobTile)
	bob.ChickenFoot = true
	bob.ChickenFootCoord = Coord{10, 8}

	// alice just extended her own line with a tile whose open end (NextPips
	// 6, at CoordA) lands adjacent to bob's exposed foot -- close the loop.
	aliceTile := &LaidTile{Tile: Tile{PipsA: 6, PipsB: 2}, Coord: Coord{11, 8}, Orientation: "right", NextPips: 6, PlayerName: "alice"}
	r.PlayerLines["alice"] = append(r.PlayerLines["alice"], aliceTile)

	r.findOuroboros(ctx, g, alice, aliceTile)

	if !bob.Dead {
		t.Error("bob should be consumed (Dead) by the ouroboros")
	}
	if bob.Score != -1 {
		t.Errorf("bob.Score = %d, want -1", bob.Score)
	}
	if !bobTile.Dead {
		t.Error("bob's line tile should be killed")
	}
	// Playing directly on your own line (lt.PlayerName != "") means you go
	// down with whoever you consume -- this is intentional (see the
	// self-referential "n == player.Name" formatting in findOuroboros'
	// note message), and is exactly the kind of easy-to-lose subtlety a
	// refactor could silently drop.
	if !alice.Dead {
		t.Error("alice should also die: consuming via her own line (not a free line) kills the consumer too")
	}
	if !aliceTile.Dead {
		t.Error("alice's own consuming tile should be killed too")
	}
	if alice.Score != 1 {
		t.Errorf("alice.Score = %d, want 1 (net: +1 for consuming bob, +1-1 wash for consuming herself)", alice.Score)
	}
	wantKills := []string{"bob", "alice"}
	if len(alice.Kills) != len(wantKills) || alice.Kills[0] != wantKills[0] || alice.Kills[1] != wantKills[1] {
		t.Errorf("alice.Kills = %v, want %v", alice.Kills, wantKills)
	}
}

func TestFindOuroborosFreeLineMoverSurvives(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	alice, bob := g.Players[0], g.Players[1]

	bobTile := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 6}, Coord: Coord{10, 9}, Orientation: "up", NextPips: 6, PlayerName: "bob"}
	r.PlayerLines["bob"] = append(r.PlayerLines["bob"], bobTile)
	bob.ChickenFoot = true
	bob.ChickenFootCoord = Coord{10, 8}

	// Same geometry as the self-consuming test, but this tile started a
	// free line (PlayerName == ""), so alice is not counted as consumed.
	freeTile := &LaidTile{Tile: Tile{PipsA: 6, PipsB: 2}, Coord: Coord{11, 8}, Orientation: "right", NextPips: 6, PlayerName: ""}
	r.FreeLines = append(r.FreeLines, []*LaidTile{freeTile})

	r.findOuroboros(ctx, g, alice, freeTile)

	if !bob.Dead {
		t.Error("bob should still be consumed")
	}
	if alice.Dead {
		t.Error("alice should survive: consuming via a free line doesn't kill the mover")
	}
	if alice.Score != 1 {
		t.Errorf("alice.Score = %d, want 1", alice.Score)
	}
	if len(alice.Kills) != 1 || alice.Kills[0] != "bob" {
		t.Errorf("alice.Kills = %v, want [bob]", alice.Kills)
	}
}

// -- killDeadLines: a line whose only escape is fully boxed in --

func TestKillDeadLinesCutoff(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice", "bob")
	alice, bob := g.Players[0], g.Players[1]

	// bob extended his line to (11,9)-(11,8); his open end is (11,8).
	bobExt := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 9}, Coord: Coord{11, 9}, Orientation: "up", NextPips: 9, PlayerName: "bob"}
	r.PlayerLines["bob"] = append(r.PlayerLines["bob"], bobExt)
	r.LaidTiles = append(r.LaidTiles, bobExt)

	// Seal every neighbor of (11,8), leaving no room for bob's line to
	// ever continue.
	for _, seal := range []Coord{{11, 7}, {10, 8}, {12, 8}} {
		r.LaidTiles = append(r.LaidTiles, &LaidTile{Tile: Tile{PipsA: 1, PipsB: 1}, Coord: seal, Orientation: "right"})
	}

	squarePips := r.MapTiles(ctx)
	r.killDeadLines(ctx, g, alice, squarePips)

	if !bob.Dead {
		t.Fatal("bob should be cut off and marked Dead")
	}
	if bob.Score != -1 {
		t.Errorf("bob.Score = %d, want -1", bob.Score)
	}
	if alice.Score != 1 {
		t.Errorf("alice.Score = %d, want 1 (credited for cutting bob off)", alice.Score)
	}
	if !bobExt.Dead {
		t.Error("bob's line tile should be marked dead")
	}
	if alice.Dead {
		t.Error("alice's own line (still just the wide-open leader) should not be cut off")
	}
	if len(alice.Kills) != 1 || alice.Kills[0] != "bob" {
		t.Errorf("alice.Kills = %v, want [bob]", alice.Kills)
	}
}

// -- Starting a free line off a spacer --

func TestLayTileStartsFreeLineOffSpacer(t *testing.T) {
	ctx := t.Context()
	g, r := newLeaderRound(4, "alice")
	alice := g.Players[0]

	// alice already extended her main line once, which is required before
	// she can branch off into a free line.
	aliceExt := &LaidTile{Tile: Tile{PipsA: 4, PipsB: 7}, Coord: Coord{9, 10}, Orientation: "left", NextPips: 7, PlayerName: "alice"}
	r.PlayerLines["alice"] = append(r.PlayerLines["alice"], aliceExt)
	r.LaidTiles = append(r.LaidTiles, aliceExt)

	r.Spacer = &Spacer{A: Coord{10, 9}, B: Coord{10, 4}}
	alice.Hand = []Tile{{PipsA: 5, PipsB: 5}, {PipsA: 9, PipsB: 9}}

	err := g.LayTile(ctx, "alice", &LaidTile{Tile: Tile{PipsA: 5, PipsB: 5}, Coord: Coord{10, 3}, Orientation: "up"})
	if err != nil {
		t.Fatalf("LayTile (free line off spacer): %v", err)
	}
	if len(r.FreeLines) != 1 {
		t.Fatalf("expected 1 free line, got %d", len(r.FreeLines))
	}
	if got := r.FreeLines[0][0].Tile; got != (Tile{PipsA: 5, PipsB: 5}) {
		t.Errorf("free line tile = %v, want 5:5", got)
	}
	if r.FreeLines[0][0].PlayerName != "" {
		t.Errorf("free line tile PlayerName = %q, want empty", r.FreeLines[0][0].PlayerName)
	}
	if r.HighestLeader != 5 {
		t.Errorf("HighestLeader = %d, want 5", r.HighestLeader)
	}
	if r.Spacer != nil {
		t.Error("the spacer should be consumed (reset to nil) once used")
	}
}
