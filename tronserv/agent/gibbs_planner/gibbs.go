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
	opportunities []*game.LaidTile
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

func (gp *GibbsPlanner) addOpportunities(ctx context.Context, g *game.Game) {
	r := g.CurrentRound(ctx)

	playerLineHeads := []*game.LaidTile{}
	for _, p := range g.Players {
		playerLineHeads = append(playerLineHeads, r.PlayerLines[p.Name][len(r.PlayerLines[p.Name])-1])
	}
	freeLineHeads := []*game.LaidTile{}
	for _, fl := range r.FreeLines {
		freeLineHeads = append(freeLineHeads, fl[len(fl)-1])
	}
	// Identify the player whose turn it is and identify the opportunities.
	nextPlayer := g.Players[g.Turn]
	nextPlayerHS := gp.hands[g.Turn]

	nextPlayerHS.opportunities = []*game.LaidTile{}
	nextPlayerHS.justDrew = false
	nextPlayerHS.justPassed = false
	nextPlayerHS.justLaid = game.Tile{}

	for pi, plh := range playerLineHeads {
		// No one can play a dead line.
		if plh.Dead {
			continue
		}
		// For other player lines...
		if pi != g.Turn {
			// You can't play them if they're not chicken footed.
			if !g.Players[pi].ChickenFoot {
				continue
			}
			// You can't play them if you're chicken footed.
			if nextPlayer.ChickenFoot {
				continue
			}
		}
		nextPlayerHS.opportunities = append(nextPlayerHS.opportunities, plh)
	}
	for _, flh := range freeLineHeads {
		// No one can play a dead line.
		if flh.Dead {
			continue
		}
		// You can't play them if you're chicken footed.
		if nextPlayer.ChickenFoot {
			continue
		}
		nextPlayerHS.opportunities = append(nextPlayerHS.opportunities, flh)
	}

	// Identify the player who went last, and fill in the action.
}

func (gp *GibbsPlanner) Update(ctx context.Context, g *game.Game) {
	if gp.lastGame == nil {
		gp.createInitialGuesses(ctx, g)
	} else {
		gp.fixBadGuesses(ctx, g)
	}

	gp.addOpportunities(ctx, g)

	//log.Printf("guessed bag: %v", gp.bag)
	for i, hs := range gp.hands {
		log.Printf("guessed hand[%d]: %v", i, hs.tiles)
		log.Printf("last opportunities[%d]: %v", i, hs.opportunities)
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
