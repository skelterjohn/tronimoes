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
		for i, p := range g.Players {
			if p.Name == gp.Name {
				gp.myPlayerIndex = i
				break
			}
		}
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
		for i, p := range g.Players {
			if i == gp.myPlayerIndex {
				// we know our own hand.
				log.Printf("our own hand: %v", p.Hand)
				gp.hands = append(gp.hands, &HandState{})
				for _, pts := range p.Hand {
					gp.RemoveTileFromBag(ctx, *pts)
					gp.hands[i].tiles = append(gp.hands[i].tiles, *pts)
				}
				continue
			}
			gp.hands = append(gp.hands, &HandState{
				tiles: gp.bag[:len(p.Hand)],
			})
			gp.bag = gp.bag[len(p.Hand):]
		}
		log.Printf("initial bag: %v", gp.bag)
		for i, hs := range gp.hands {
			log.Printf("initial hand[%d]: %v", i, hs.tiles)
		}
	}

	// Ensure any laid tiles are removed from our guessed bag, and
	// if it's in a hand, swap it with one in the bag to keep the hand
	// size correct (and guess at what they may have instead).
	laidTiles := g.CurrentRound(ctx).LaidTiles
	for _, lt := range laidTiles {
		if gp.RemoveTileFromBag(ctx, *lt.Tile) {
			continue
		}
		for i := range gp.hands {
			if gp.RemoveTileFromHand(ctx, i, *lt.Tile) {
				break
			}
		}
	}

	for i, p := range g.Players {
		if i == gp.myPlayerIndex {
			// for our own hand, check the guessed bag and other hands. If
			// we remove it from one another hand, replace it with one from
			// the bag.
			for _, ht := range gp.hands[i].tiles {
				if gp.RemoveTileFromBag(ctx, ht) {
					continue
				}
				for oi := range gp.hands {
					if oi == i {
						continue
					}
					if gp.RemoveTileFromHand(ctx, oi, ht) {
						gp.hands[oi].tiles = append(gp.hands[oi].tiles, gp.bag[0])
						gp.bag = gp.bag[1:]
						break
					}
				}
			}
			continue
		}
		extraTiles := len(p.Hand) - len(gp.hands[i].tiles)
		if extraTiles > 0 {
			// add tiles from the bag (they must have drawn)
			gp.hands[i].tiles = append(gp.hands[i].tiles, gp.bag[:extraTiles]...)
			gp.bag = gp.bag[extraTiles:]
		}
		if extraTiles < 0 {
			// remove tiles from the hand (they must have laid)
			removeCount := -extraTiles
			removedTiles := gp.hands[i].tiles[:removeCount]
			gp.bag = append(gp.bag, removedTiles...)
			gp.hands[i].tiles = gp.hands[i].tiles[removeCount:]
		}
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
