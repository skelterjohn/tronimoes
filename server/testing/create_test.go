package testing

import (
	"context"
	"testing"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"google.golang.org/grpc/metadata"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	c, close := createBufferedServer(t, ctx)
	defer close()

	gameID := gameForPlayers(t, ctx, c, []string{"jt", "stef"})

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
