package gibbs_planner

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/skelterjohn/tronimoes/tronserv/agent/reacts"
	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/client"
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
	Name                  string
	MaxInferenceTime      time.Duration
	MaxSimulationTime     time.Duration
	MaxSimulationDepth    int
	MaxSimulationsPerMove int
	ValueDecay            float64

	lastGame      *game.Game
	bag           []game.Tile
	hands         []*HandState
	myPlayerIndex int
	Client        *client.TronimoesClient
	currentRoot   *PlanNode
}

func (gp *GibbsPlanner) SetDefaults() {
	gp.MaxInferenceTime = 1 * time.Second
	gp.MaxSimulationTime = 5 * time.Second
	gp.MaxSimulationDepth = 15
	gp.ValueDecay = 0.9
	gp.MaxSimulationsPerMove = 0
}

func (gp *GibbsPlanner) Ready(ctx context.Context) {
	gp.React(ctx, "bow")
	gp.currentRoot = nil
}

func (gp *GibbsPlanner) CheckScore(ctx context.Context, pg, g *game.Game) {
	if pg != nil {
		scoreDiff := g.GetPlayer(ctx, gp.Name).Score - pg.GetPlayer(ctx, gp.Name).Score
		if scoreDiff > 0 {
			gp.React(ctx, "happy")
		} else if scoreDiff < 0 {
			gp.React(ctx, "sad")
		}
	}
}

func (gp *GibbsPlanner) Update(ctx context.Context, previousGame *game.Game, g *game.Game) {
	gp.CheckScore(ctx, previousGame, g)

	if gp.lastGame == nil || len(g.Rounds) != len(previousGame.Rounds) {
		gp.createInitialGuesses(ctx, g)
	} else {
		gp.fixBadGuesses(ctx, g)
	}

	gp.addOpportunities(ctx, previousGame, g)
	gp.lastGame = g

	ctx, cancel := context.WithTimeout(ctx, gp.MaxInferenceTime)
	defer cancel()
	gp.ConsiderSwaps(ctx)
	for i := range gp.hands {
		if i == gp.myPlayerIndex {
			continue
		}
		game.Debug(ctx, "sampled hand[%d]: %v", i, gp.hands[i].tiles)
	}
}

func (gp *GibbsPlanner) GetMove(ctx context.Context, g *game.Game, p *game.Player) types.Move {
	gp.CheckScore(ctx, gp.lastGame, g)
	legalMoves, legalSpacers := g.CurrentRound(ctx).FindLegalMoves(ctx, g, p)

	playingOffRoundLeader := len(g.CurrentRound(ctx).PlayerLines[p.Name]) == 1
	// If it's the round leader, we still have a choice to make that can be
	// improved with planning. Otherwise, just pass/draw without planning.
	if !playingOffRoundLeader && len(legalMoves) == 0 && len(legalSpacers) == 0 {
		if p.JustDrew || len(g.Bag) == 0 {
			game.Log(ctx, "no legal moves or spacers, and just drew; passing")
			return types.Move{
				Pass:     true,
				Selected: types.RandomInitialFoot(g),
			}
		}
		game.Log(ctx, "no legal moves or spacers, and haven't drawn yet; drawing")
		return types.Move{
			Draw: true,
		}
	}

	root := NewPlanNode(g.Turn, len(gp.hands), 0)

	ctx, cancel := context.WithTimeout(ctx, gp.MaxSimulationTime)
	defer cancel()

	for i, hs := range gp.hands {
		if i == gp.myPlayerIndex {
			continue
		}
		g.Players[i].Hand = nil
		for _, t := range hs.tiles {
			g.Players[i].Hand = append(g.Players[i].Hand, t)
		}
	}
	g.Bag = nil
	for _, t := range gp.bag {
		g.Bag = append(g.Bag, t)
	}

	gdata, err := json.Marshal(g)
	if err != nil {
		game.Debug(ctx, "error marshalling game: %v", err)
	}

	simulating := true
	simulations := 0
	for simulating {
		select {
		case <-ctx.Done():
			simulating = false
		default:
		}
		var sg game.Game
		if err := json.Unmarshal(gdata, &sg); err != nil {
			game.Debug(ctx, "error unmarshalling game: %v", err)
		}
		if err := gp.SimulateGame(ctx, &sg, root, gp.MaxSimulationDepth); err != nil {
			game.Debug(ctx, "error simulating game: %v", err)
		}
		simulations++
		if gp.MaxSimulationsPerMove > 0 && simulations >= gp.MaxSimulationsPerMove {
			simulating = false
		}
	}
	game.Log(ctx, "simulated %d games", simulations)
	bestMove := root.ChooseBestMove(ctx)
	game.Debug(ctx, "hand: %v", g.Players[g.Turn].Hand)
	game.Debug(ctx, "best move: %s %v", bestMove, root.Moves[bestMove].V)

	if bestMove.Pass {
		gp.React(ctx, "frustration")
	}
	if bestMove.PlaceSpacer {
		gp.React(ctx, "free")
	}
	return bestMove
}

func (gp *GibbsPlanner) CompleteRound(ctx context.Context, g *game.Game) {
	p := g.GetPlayer(ctx, gp.Name)
	if p.Dead {
		gp.React(ctx, "skull")
	}
	if len(p.Hand) == 0 {
		gp.React(ctx, "work")
	}
	othersLive := false
	for _, op := range g.Players {
		if op.Name == gp.Name {
			continue
		}
		if !op.Dead {
			othersLive = true
		}
	}
	if !othersLive {
		gp.React(ctx, "zap")
	}
}

func (gp *GibbsPlanner) CompleteGame(ctx context.Context, g *game.Game) {
	p := g.GetPlayer(ctx, gp.Name)
	highScore := p.Score
	for _, op := range g.Players {
		if op.Score > highScore {
			highScore = op.Score
		}
	}
	if p.Score == highScore {
		gp.ReactWait(ctx, "victory")
	}
}

func (gp *GibbsPlanner) React(ctx context.Context, query string) {
	go func(ctx context.Context) {
		gp.ReactWait(ctx, query)
	}(context.WithoutCancel(ctx))
}

func (gp *GibbsPlanner) ReactWait(ctx context.Context, query string) {
	game.Log(ctx, "reacting: %s", query)
	url, err := reacts.FindImageURL(ctx, query)
	if err != nil {
		game.Log(ctx, "Error getting image URL: %v", err)
		return
	}
	if _, err := gp.Client.React(ctx, url); err != nil {
		game.Log(ctx, "Error reacting: %v", err)
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
					if pj <= pi {
						continue
					}
					gp.ConsiderSwapHands(ctx, pi, pj)
				}
			}
		}
	}
}

func (gp *GibbsPlanner) ConsiderSwapHands(ctx context.Context, pi, pj int) {
	for i, ti := range gp.hands[pi].tiles {
		scores := make([]float64, len(gp.hands[pj].tiles)+1)
		for j, tj := range gp.hands[pj].tiles {
			scores[j] = gp.ScoreInHand(ctx, gp.hands[pi], tj, ti)
		}
		scores[len(gp.hands[pj].tiles)] = gp.ScoreInHand(ctx, gp.hands[pi], ti, ti)
		chosenIndex := ChooseIndex(scores)
		if chosenIndex == len(gp.hands[pj].tiles) {
			continue
		}
		gp.hands[pi].tiles[i], gp.hands[pj].tiles[chosenIndex] = gp.hands[pj].tiles[chosenIndex], gp.hands[pi].tiles[i]
	}
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
