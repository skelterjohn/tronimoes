package testing

import (
	"context"
	"sync"
	"testing"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"google.golang.org/grpc/metadata"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	c, close := createBufferedServer(t, ctx)
	defer close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var g1, g2 *spb.Game
	var err1, err2 error
	go func() {
		defer wg.Done()
		g1, err1 = createGameAndWait(t, ctx, c, "jt", &spb.CreateGameRequest{
			Discoverable: false,
			Private:      false,
			MinPlayers:   0,
			MaxPlayers:   0,
			BoardShape:   spb.CreateGameRequest_standard_31_by_30,
		})
	}()
	go func() {
		defer wg.Done()
		g2, err2 = createGameAndWait(t, ctx, c, "stef", &spb.CreateGameRequest{
			Discoverable: false,
			Private:      false,
			MinPlayers:   0,
			MaxPlayers:   0,
			BoardShape:   spb.CreateGameRequest_standard_31_by_30,
		})
	}()
	wg.Wait()
	if err1 != nil {
		t.Fatalf("Could not create game 1: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("Could not create game 2: %v", err2)
	}

	if g1.GetGameId() != g2.GetGameId() {
		t.Errorf("Game IDs did not match, %q != %q", g1.GetGameId(), g2.GetGameId())
	}

	gameID := g1.GetGameId()

	checkBoardForPlayer := func(t *testing.T, ctx context.Context, playerID string) {
		t.Helper()

		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			"access_token": playerID,
			"player_id":    playerID,
		}))

		b, err := c.GetBoard(ctx, &spb.GetBoardRequest{
			GameId: gameID,
		})
		if err != nil {
			t.Errorf("Could not get board for %s: %v", playerID, err)
			return
		}
		for _, tile := range b.GetBag() {
			if tile.A != -1 {
				t.Errorf("%s can see the bag contents: %v", playerID, tile)
				break
			}
		}
		for _, p := range b.GetPlayers() {
			for _, tile := range p.GetHand() {
				if p.GetPlayerId() == playerID {
					if tile.A == -1 {
						t.Errorf("%s cannot see own hand: %v", playerID, tile)
						break
					}
					continue
				}
				if tile.A != -1 {
					t.Errorf("%s can see the %s's hand: %v", playerID, p.GetPlayerId(), tile)
					break
				}
			}
		}
	}
	checkBoardForPlayer(t, ctx, "stef")
	checkBoardForPlayer(t, ctx, "jt")
}
