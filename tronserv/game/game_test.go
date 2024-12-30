package game

import (
	"context"
	"testing"
)

func TestDrawForLeader(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 30; i++ {
		game := NewGame(ctx, "ABCDEF")

		game.AddPlayer(ctx, &Player{
			Name:  "p1",
			Ready: true,
		})
		game.AddPlayer(ctx, &Player{
			Name:  "p2",
			Ready: true,
		})

		prevRound := &Round{
			LaidTiles: []*LaidTile{
				{Tile: &Tile{PipsA: 1, PipsB: 1}},
			},
			Done: true,
		}
		game.Rounds = append(game.Rounds, prevRound)

		game.Start(ctx, "p1")

		if game.CheckForDupes(ctx, "test") {
			t.Errorf("dupes found")
		}
	}
}
