package main

import (
	"context"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type RandomChoice struct {
}

func (RandomChoice) Ready(ctx context.Context) {

}
func (RandomChoice) Update(ctx context.Context, g *game.Game) {

}
func (RandomChoice) GetMove(ctx context.Context, g *game.Game, p *game.Player) Move {
	legalMoves, legalSpacers := g.CurrentRound(ctx).FindLegalMoves(ctx, g, p)

	if len(legalSpacers) > 0 {
		return Move{
			Spacer: legalSpacers[rand.Intn(len(legalSpacers))],
		}
	}
	if len(legalMoves) > 0 {
		return Move{
			LaidTile: legalMoves[rand.Intn(len(legalMoves))],
		}
	}
	if p.JustDrew {
		return Move{
			Pass: true,
			Selected: game.Coord{
				X: g.BoardWidth / 2,
				Y: (g.BoardHeight / 2) - 1,
			},
		}
	}
	return Move{
		Draw: true,
	}
}
