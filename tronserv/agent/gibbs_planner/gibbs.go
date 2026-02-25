package gibbs_planner

import (
	"context"
	"log"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type HandState struct {
	tiles         []game.Tile
	justDrew      bool
	justPassed    bool
	justLaid      game.Tile
	opportunities []game.Tile
}

type GibbsPlanner struct {
	Name          string
	lastGame      *game.Game
	bag           []game.Tile
	hands         []*HandState
	myPlayerIndex int
}

func (GibbsPlanner) Ready(ctx context.Context) {

}

func (gp *GibbsPlanner) RemoveTileFromBag(ctx context.Context, tile game.Tile) bool {
	for i := range gp.bag {
		if gp.bag[i] == tile {
			gp.bag[i] = gp.bag[len(gp.bag)-1]
			gp.bag = gp.bag[:len(gp.bag)-1]
			return true
		}
	}
	return false
}

func (gp *GibbsPlanner) RemoveTileFromHand(ctx context.Context, whichPlayer int, tile game.Tile) bool {
	hs := gp.hands[whichPlayer]
	for i, ht := range hs.tiles {
		if ht == tile {
			hs.tiles[i] = hs.tiles[len(hs.tiles)-1]
			hs.tiles = hs.tiles[:len(hs.tiles)-1]
			return true
		}
	}
	return false
}

func (gp *GibbsPlanner) Update(ctx context.Context, g *game.Game) {
	if gp.lastGame == nil {
		gp.createInitialGuesses(ctx, g)
	} else {
		gp.fixBadGuesses(ctx, g)
	}

	//log.Printf("guessed bag: %v", gp.bag)
	for i, hs := range gp.hands {
		log.Printf("guessed hand[%d]: %v", i, hs.tiles)
	}

	gp.lastGame = g
}
func (gp *GibbsPlanner) GetMove(ctx context.Context, g *game.Game, p *game.Player) types.Move {
	gp.Update(ctx, g)
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
