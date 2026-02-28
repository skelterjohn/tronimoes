package main

import (
	"context"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type RandomChoice struct {
}

func (RandomChoice) Ready(ctx context.Context) {

}
func (RandomChoice) Update(ctx context.Context, previousGame *game.Game, g *game.Game) {

}
func (RandomChoice) GetMove(ctx context.Context, g *game.Game, p *game.Player) types.Move {
	legalMoves, legalSpacers := g.CurrentRound(ctx).FindLegalMoves(ctx, g, p)

	if len(legalSpacers) > 0 {
		return types.Move{
			Spacer: legalSpacers[rand.Intn(len(legalSpacers))],
		}
	}
	if len(legalMoves) > 0 {
		return types.Move{
			LaidTile: legalMoves[rand.Intn(len(legalMoves))],
		}
	}
	if p.JustDrew {
		// randomly choose one, so if it's bad we randomly choose again.
		cfSelection := game.Coord{
			X: g.BoardWidth / 2,
			Y: (g.BoardHeight / 2),
		}
		var dx = rand.Intn(2)
		dy := rand.Intn(3) - 1
		if dy == 0 {
			if dx == 0 {
				dx = -1
			} else {
				dx = 2
			}
		}
		cfSelection.X += dx
		cfSelection.Y += dy

		return types.Move{
			Pass:     true,
			Selected: cfSelection,
		}
	}
	return types.Move{
		Draw: true,
	}
}
