package main

import (
	"context"
	"io"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type RandomChoice struct {
}

func (RandomChoice) Configure(ctx context.Context, r io.Reader) error {
	return nil
}

func (RandomChoice) Ready(ctx context.Context) {

}
func (RandomChoice) Update(ctx context.Context, previousGame *game.Game, g *game.Game) {

}
func (RandomChoice) GetMove(ctx context.Context, g *game.Game, p *game.Player) types.Move {
	legalMoves, legalSpacers, legalPassFeet := g.CurrentRound(ctx).FindLegalMoves(ctx, g, p)

	if len(legalSpacers) > 0 {
		return types.Move{
			PlaceSpacer: true,
			Spacer:      legalSpacers[rand.Intn(len(legalSpacers))],
		}
	}
	if len(legalMoves) > 0 {
		return types.Move{
			LayTile:  true,
			LaidTile: legalMoves[rand.Intn(len(legalMoves))],
		}
	}
	if p.JustDrew || len(g.Bag) == 0 {
		var selected game.Coord
		if len(legalPassFeet) > 0 {
			selected = legalPassFeet[rand.Intn(len(legalPassFeet))]
		}
		// randomly choose one, so if it's bad we randomly choose again.
		return types.Move{
			Pass:     true,
			Selected: selected,
		}
	}
	return types.Move{
		Draw: true,
	}
}

func (RandomChoice) CompleteRound(ctx context.Context, g *game.Game) {

}

func (RandomChoice) CompleteGame(ctx context.Context, g *game.Game) {

}
