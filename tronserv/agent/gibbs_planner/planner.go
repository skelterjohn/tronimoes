package gibbs_planner

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type PlanNode struct {
	Depth   int
	Visited int
	Turn    int
	Moves   map[string]*PlanNode
	Eval    []float64
	R       []float64
	V       []float64
	H       []float64
}

func NewPlanNode(turn int, count int, depth int) *PlanNode {
	return &PlanNode{
		Turn:  turn,
		Depth: depth,
		Moves: make(map[string]*PlanNode),
		Eval:  make([]float64, count),
		R:     make([]float64, count),
		V:     make([]float64, count),
		H:     make([]float64, count),
	}
}

func (n *PlanNode) Next(move types.Move, turn, count int) *PlanNode {
	// fmt.Printf("Next %s\n", move)
	moveStr := move.JSON()
	nextNode, ok := n.Moves[moveStr]
	if !ok {
		nextNode = NewPlanNode(turn, count, n.Depth+1)
		n.Moves[moveStr] = nextNode
	}
	return nextNode
}

func (n *PlanNode) ChooseBestMove(ctx context.Context) (types.Move, error) {
	bestMoveStr := ""
	bestV := -math.MaxFloat64
	for moveStr, next := range n.Moves {
		nextV := next.V[n.Turn]
		if nextV > bestV {
			bestV = nextV
			bestMoveStr = moveStr
		}
	}
	return types.MoveFromJSON(bestMoveStr)
}

// We need to be given a fresh game copy, because it's gonna get messed up.
func (gp *GibbsPlanner) SimulateGame(ctx context.Context, g *game.Game, root *PlanNode, maxDepth int) error {
	// Play a random game until it's done or we reach max depth or it takes too long.
	curNode := root
	nodesInSimulation := []*PlanNode{root}
	r := g.CurrentRound(ctx)

	for i := range root.Eval {
		root.Eval[i] = float64(g.Players[i].Score)
	}
	game.Debug(ctx, "Simulating game at depth %d", maxDepth)

	for !r.Done && maxDepth > 0 {
		select {
		case <-ctx.Done():
			maxDepth = 0
		default:
			maxDepth--
		}
		legalMoves, legalSpacers := r.FindLegalMoves(ctx, g, g.Players[g.Turn])

		game.Debug(ctx, "%s has %d tiles, %d spacers", g.Players[g.Turn].Name, len(legalMoves), len(legalSpacers))
		moveCount := len(legalMoves) + len(legalSpacers)
		moveCount += 1 // draw or pass

		unnormalizedLogLikelihoods := make([]float64, moveCount)
		for i := range unnormalizedLogLikelihoods {
			unnormalizedLogLikelihoods[i] = 0
		}

		for i, lt := range legalMoves {
			lt.PlayerName = g.Players[g.Turn].Name
			nn := curNode.Next(types.Move{LaidTile: lt}, g.Turn, len(gp.hands))
			// bias the planner towards high-value, away from low-value.
			unnormalizedLogLikelihoods[i] += nn.V[g.Turn]
			// bias the planner away from options that have been considered a lot.
			if nn.Visited > 1000 {
				unnormalizedLogLikelihoods[i] -= 1000
			} else if nn.Visited > 100 {
				unnormalizedLogLikelihoods[i] -= 5
			} else if nn.Visited > 10 {
				unnormalizedLogLikelihoods[i] -= 1
			}
		}
		for i := range legalSpacers {
			// free lines are exciting, and good opportunities to harass opponents.
			unnormalizedLogLikelihoods[i+len(legalMoves)] += 1
		}
		// drawing is less exciting, so let's examine it less than other options.
		// the agent will still draw when it has to.
		whichMove := ChooseIndex(unnormalizedLogLikelihoods)

		var nextNode *PlanNode
		var bestMove types.Move

		if whichMove < len(legalMoves) {
			move := legalMoves[whichMove]
			game.Debug(ctx, "p%d lays %s", g.Turn, move)
			if err := g.LayTile(ctx, g.Players[g.Turn].Name, move); err != nil {
				return fmt.Errorf("laying: %w", err)
			}
			bestMove = types.Move{LaidTile: move}
		} else if whichMove == moveCount-1 {
			if !g.Players[g.Turn].JustDrew {
				if !g.DrawTile(ctx, g.Players[g.Turn].Name) {
					return errors.New("drawing failed")
				}
				bestMove = types.Move{Draw: true}
			} else {
				cfSelection := types.RandomInitialFoot(g)
				if err := g.Pass(ctx, g.Players[g.Turn].Name, cfSelection.X, cfSelection.Y); err != nil {
					return fmt.Errorf("passing: %w", err)
				}
				bestMove = types.Move{Pass: true, Selected: cfSelection}
			}
		} else {
			spacer := legalSpacers[whichMove-len(legalMoves)]
			if err := g.LaySpacer(ctx, g.Players[g.Turn].Name, spacer); err != nil {
				return fmt.Errorf("spacing: %w", err)
			}
			bestMove = types.Move{Spacer: spacer}
		}

		game.Debug(ctx, "p%d -> %s", curNode.Turn, bestMove)

		nextNode = curNode.Next(bestMove, g.Turn, len(gp.hands))
		nextNode.Eval = make([]float64, len(gp.hands))
		for i := range nextNode.Eval {
			nextNode.Eval[i] = float64(g.Players[i].Score)
		}
		nextNode.R = make([]float64, len(gp.hands))
		for i := range nextNode.R {
			nextNode.R[i] = nextNode.Eval[i] - curNode.Eval[i]
		}
		game.Debug(ctx, " R <- %v from %v-%v", nextNode.R, nextNode.Eval, curNode.Eval)
		nodesInSimulation = append(nodesInSimulation, nextNode)
		curNode = nextNode
	}

	if !r.Done {
		curNode.H = gp.Heuristic(ctx, g, root)
		game.Debug(ctx, "Heuristic: %v @ %d", curNode.H, curNode.Depth)
	}
	// The rest is fast so we still do it if we ran out of time.

	// Now that we've simulated a random game, let's backprop its eval.

	// First we start with the score at the end of this simulation.

	// lastNode := nodesInSimulation[len(nodesInSimulation)-1]
	curNode.Eval = make([]float64, len(gp.hands))
	for i := range curNode.Eval {
		curNode.Eval[i] = float64(g.Players[i].Score)
	}
	// fmt.Printf("last V: %v\n", curNode.V)
	// fmt.Printf("last R: %v\n", curNode.R)

	// For every node, we assume that the player is maximizing their own score.
	// This is myopic to the round, missing outcomes that might win the game. Alas.
	for i := len(nodesInSimulation) - 1; i >= 0; i-- {
		cur := nodesInSimulation[i]
		bestMove, err := cur.ChooseBestMove(ctx)
		if err != nil {
			return fmt.Errorf("choosing best move: %w", err)
		}
		bestNode := cur.Moves[bestMove.JSON()]
		if bestNode == nil {
			copy(cur.V, cur.R)
			continue
		}
		game.Debug(ctx, "best from %d +1 H: %v V: %v", cur.Depth, bestNode.H, bestNode.V)
		for i, bv := range bestNode.V {
			vh := bestNode.H[i] + bv
			cur.V[i] = gp.ValueDecay*vh + cur.R[i]
		}
	}

	for i, n := range nodesInSimulation {
		n.Visited++
		game.Debug(ctx, "%d: p%d @ %d / %d", i, n.Turn, n.Depth, n.Visited)
		game.Debug(ctx, "   V: %v", n.V)
		game.Debug(ctx, "   R: %v", n.R)
		game.Debug(ctx, "   H: %v", n.H)
	}
	return nil
}
