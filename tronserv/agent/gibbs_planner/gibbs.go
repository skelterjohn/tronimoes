package gibbs_planner

import (
	"context"
	"log"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type HandState struct {
	tiles []game.Tile
}

type GibbsPlanner struct {
	lastGame *game.Game
	bag      []game.Tile
	hands    []HandState
}

func (GibbsPlanner) Ready(ctx context.Context) {

}
func (gp *GibbsPlanner) Update(ctx context.Context, g *game.Game) {
	if gp.lastGame == nil {
		// A round just started, let's halucinate a bag and deal it out.
		gp.bag = nil
		for a := 0; a <= g.MaxPips; a++ {
			for b := a; b <= g.MaxPips; b++ {
				gp.bag = append(gp.bag, game.Tile{
					PipsA: a,
					PipsB: b,
				})
			}
		}
		rand.Shuffle(len(gp.bag), func(i, j int) {
			gp.bag[i], gp.bag[j] = gp.bag[j], gp.bag[i]
		})
		for _, p := range g.Players {
			gp.hands = append(gp.hands, HandState{
				tiles: gp.bag[:len(p.Hand)],
			})
			gp.bag = gp.bag[len(p.Hand):]
		}
	}

	// Go through all the laid tiles and ensure they're removed.
tileLoop:
	for _, lt := range g.CurrentRound(ctx).LaidTiles {
		// first get them out of the bag
		for i := range gp.bag {
			if gp.bag[i] == *lt.Tile {
				gp.bag[i] = gp.bag[len(gp.bag)-1]
				gp.bag = gp.bag[:len(gp.bag)-1]
				continue tileLoop
			}
		}
		// then get them out of the hands
		for _, hs := range gp.hands {
			for i := range hs.tiles {
				if hs.tiles[i] == *lt.Tile {
					hs.tiles[i] = gp.bag[len(gp.bag)-1]
					gp.bag = gp.bag[:len(gp.bag)-1]
					continue tileLoop
				}
			}
		}
	}

	log.Printf("guessed bag: %v", gp.bag)
	for i, hs := range gp.hands {
		log.Printf("guessed hand for %d: %v", i, hs.tiles)
	}

	gp.lastGame = g
}
func (GibbsPlanner) GetMove(ctx context.Context, g *game.Game, p *game.Player) types.Move {
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
		return types.Move{
			Pass: true,
			Selected: game.Coord{
				X: g.BoardWidth / 2,
				Y: (g.BoardHeight / 2) - 1,
			},
		}
	}
	return types.Move{
		Draw: true,
	}
}
