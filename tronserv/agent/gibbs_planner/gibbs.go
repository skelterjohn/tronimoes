package gibbs_planner

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type HandState struct {
	tiles         []game.Tile
	justDrew      bool
	justPassed    bool
	justLaid      *game.LaidTile
	opportunities []*game.LaidTile
}

func (hs *HandState) String() string {
	return fmt.Sprintf("HandState{tiles: %v, justDrew: %v, justPassed: %v, justLaid: %s, opportunities: %s}", hs.tiles, hs.justDrew, hs.justPassed, hs.justLaid, hs.opportunities)
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

func (gp *GibbsPlanner) Update(ctx context.Context, previousGame *game.Game, g *game.Game) {
	if gp.lastGame == nil || len(g.Rounds) != len(previousGame.Rounds) {
		gp.createInitialGuesses(ctx, g)
	} else {
		gp.fixBadGuesses(ctx, g)
	}

	gp.addOpportunities(ctx, previousGame, g)
	gp.lastGame = g

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	gp.ConsiderSwaps(ctx)
	for i := range gp.hands {
		if i == gp.myPlayerIndex {
			continue
		}
		log.Printf("sampled hand[%d]: %v", i, gp.hands[i].tiles)
	}
}

func (gp *GibbsPlanner) GetMove(ctx context.Context, g *game.Game, p *game.Player) types.Move {
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

func (gp *GibbsPlanner) ConsiderSwaps(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for pi := range gp.hands {
				if pi == gp.myPlayerIndex {
					continue
				}
				gp.ConsiderSwapBag(ctx, pi)
				for pj := range gp.hands {
					if pj == gp.myPlayerIndex {
						continue
					}
					if pj == pi {
						continue
					}
					gp.ConsiderSwapHands(ctx, pi, pj)
				}
			}
		}
	}
}

func (gp *GibbsPlanner) ConsiderSwapHands(ctx context.Context, pi, pj int) {

}

func (gp *GibbsPlanner) ConsiderSwapBag(ctx context.Context, pi int) {
	scores := make([]float64, len(gp.bag)+1)
	for i, ti := range gp.hands[pi].tiles {
		for j, tj := range gp.bag {
			scores[j] = gp.ScoreInHand(ctx, gp.hands[pi], tj, ti)
			// log.Printf("score %s->%s: %f", ti, tj, scores[j])
		}
		// don't forget the score of this tile staying put.
		scores[len(gp.bag)] = gp.ScoreInHand(ctx, gp.hands[pi], ti, ti)
		// log.Printf("score %s stays: %f", ti, scores[len(gp.bag)])
		chosenIndex := ChooseIndex(scores)
		if chosenIndex == len(gp.bag) {
			// log.Printf("swap[%d]: %s stays", pi, ti)
			continue
		}
		// log.Printf("swap[%d]: %s -> %s", pi, ti, gp.bag[chosenIndex])
		gp.hands[pi].tiles[i], gp.bag[chosenIndex] = gp.bag[chosenIndex], gp.hands[pi].tiles[i]
	}
}

func (gp *GibbsPlanner) ScoreInHand(ctx context.Context, hs *HandState, t, ignoredTile game.Tile) float64 {
	score := 0.0
	for _, ti := range hs.tiles {
		if ti == ignoredTile || ti == t {
			continue
		}
		canConnect := ti.PipsA == t.PipsA || ti.PipsB == t.PipsB || ti.PipsA == t.PipsB || ti.PipsB == t.PipsA
		if !canConnect {
			score -= .1
		}
	}
	if len(hs.opportunities) > 0 {
		for _, opp := range hs.opportunities {
			if t.PipsA == opp.NextPips || t.PipsB == opp.NextPips {
				score -= 1
			}
		}
	}
	return score
}

// ChooseIndex picks an index from unnormalizedLogLikelihoods by normalizing
// them, converting to likelihoods (exp), and making a weighted random choice.
func ChooseIndex(unnormalizedLogLikelihoods []float64) int {
	if len(unnormalizedLogLikelihoods) == 0 {
		return 0
	}
	// Subtract max for numerical stability (avoids exp overflow; max log becomes 0)
	maxLog := unnormalizedLogLikelihoods[0]
	for _, logL := range unnormalizedLogLikelihoods[1:] {
		if logL > maxLog {
			maxLog = logL
		}
	}
	// Convert to likelihoods: exp(logL - maxLog)
	likelihoods := make([]float64, len(unnormalizedLogLikelihoods))
	var total float64
	for i, logL := range unnormalizedLogLikelihoods {
		likelihoods[i] = math.Exp(logL - maxLog)
		// log.Printf("%d: %f", i, likelihoods[i])
		total += likelihoods[i]
	}
	// log.Printf("total: %f", total)
	if total == 0 {
		// All were -inf or underflow; choose uniformly
		return rand.Intn(len(unnormalizedLogLikelihoods))
	}
	// Weighted choice: roll in [0, total)
	r := rand.Float64() * total
	// log.Printf("r: %f", r)
	for i, w := range likelihoods {
		r -= w
		if r < 0 {
			return i
		}
	}
	return len(unnormalizedLogLikelihoods) - 1
}
