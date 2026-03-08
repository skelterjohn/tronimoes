package gibbs_planner

import (
	"context"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

func (gp *GibbsPlanner) createInitialGuesses(ctx context.Context, g *game.Game) {
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

	for _, t := range g.CurrentRound(ctx).LaidTiles {
		gp.RemoveTileFromBag(ctx, *t.Tile)
	}

	gp.hands = make([]*HandState, len(g.Players))
	for i, p := range g.Players {
		if i != gp.myPlayerIndex {
			continue
		}
		// we know our own hand.
		gp.hands[i] = &HandState{}
		for _, pts := range p.Hand {
			gp.RemoveTileFromBag(ctx, *pts)
			gp.hands[i].tiles = append(gp.hands[i].tiles, *pts)
		}
		break
	}
	for i, p := range g.Players {
		if i == gp.myPlayerIndex {
			continue
		}
		gp.hands[i] = &HandState{
			tiles: gp.bag[:len(p.Hand)],
		}
		gp.bag = gp.bag[len(p.Hand):]
	}
	Log(ctx, "initial bag (%d): %v", len(gp.bag), gp.bag)
	for i, hs := range gp.hands {
		Log(ctx, "initial hand[%d]: %v", i, hs.tiles)
	}
}

func (gp *GibbsPlanner) fixBadGuesses(ctx context.Context, g *game.Game) {

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
}

func (gp *GibbsPlanner) RemoveTileFromBag(ctx context.Context, tile game.Tile) bool {
	for i := range gp.bag {
		if gp.bag[i] == tile {
			gp.bag[i] = gp.bag[0]
			gp.bag = gp.bag[1:]
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

func (gp *GibbsPlanner) addOpportunities(ctx context.Context, previousGame, g *game.Game) {
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
	nextPlayerHS.justLaid = nil

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
	nextPlayerHS.justDrew = false
	nextPlayerHS.justPassed = false
	nextPlayerHS.justLaid = nil

	// Identify the player who went last, and fill in the action.
	pr := previousGame.CurrentRound(ctx)
	if pr != nil && (previousGame.Turn != g.Turn || len(r.LaidTiles) > len(pr.LaidTiles)) {
		pi := previousGame.Turn
		lastPlayerHS := gp.hands[pi]
		lastPlayerHS.justPassed = g.Players[pi].ChickenFoot
		lastPlayerHS.justDrew = lastPlayerHS.justPassed || len(g.Players[pi].Hand) > len(previousGame.Players[pi].Hand)
		if len(r.LaidTiles) > len(pr.LaidTiles) {
			lastPlayerHS.justLaid = r.LaidTiles[len(r.LaidTiles)-1]
		}
		Log(ctx, "inference[%d]: %s", pi, lastPlayerHS)
		Log(ctx, "In the bag: %d", len(gp.bag))
	}
}
