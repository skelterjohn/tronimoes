package gibbs_planner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type PlanNode struct {
	Turn  int
	Moves map[string]*PlanNode
	Eval  []float64
}

func NewPlanNode(turn int, count int) *PlanNode {
	return &PlanNode{
		Turn:  turn,
		Moves: make(map[string]*PlanNode),
		Eval:  make([]float64, count),
	}
}

func (n *PlanNode) Next(move string, turn, count int) *PlanNode {
	// log.Printf("Next %s", move)
	nextNode, ok := n.Moves[move]
	if !ok {
		nextNode = NewPlanNode(turn, count)
		n.Moves[move] = nextNode
	}
	return nextNode
}

func (n *PlanNode) ChooseBestMove(ctx context.Context) (string, error) {
	bestMove := ""
	bestEval := -math.MaxFloat64
	for move, next := range n.Moves {
		nextEval := next.Eval[n.Turn]
		if nextEval > bestEval {
			bestEval = nextEval
			bestMove = move
		}
	}
	return bestMove, nil
}

// We need to be given a fresh game copy, because it's gonna get messed up.
func (gp *GibbsPlanner) SimulateGame(ctx context.Context, g *game.Game, root *PlanNode, maxDepth int) error {
	// Play a random game until it's done or we reach max depth or it takes too long.
	curNode := root
	nodesInSimulation := []*PlanNode{root}
	r := g.CurrentRound(ctx)

	for !r.Done && maxDepth > 0 {
		select {
		case <-ctx.Done():
			maxDepth = 0
		default:
			maxDepth--
		}
		legalMoves, legalSpacers := r.FindLegalMoves(ctx, g, g.Players[g.Turn])
		for _, m := range legalMoves {
			m.NextPips = -1
		}
		// log.Printf("%s has %d tiles, %d spacers", g.Players[g.Turn].Name, len(legalMoves), len(legalSpacers))
		moveCount := len(legalMoves) + len(legalSpacers)
		moveCount += 1 // draw or pass
		whichMove := rand.Intn(moveCount)
		// log.Printf("whichMove: %d", whichMove)
		if whichMove < len(legalMoves) {
			move := legalMoves[whichMove]
			move.PlayerName = g.Players[g.Turn].Name
			if err := g.LayTile(ctx, g.Players[g.Turn].Name, move); err != nil {
				return fmt.Errorf("laying: %w", err)
			}
			move.NextPips = -1
			curNode = curNode.Next(move.String(), g.Turn, len(gp.hands))
			nodesInSimulation = append(nodesInSimulation, curNode)
		} else if whichMove == moveCount-1 {
			if !g.Players[g.Turn].JustDrew {
				if !g.DrawTile(ctx, g.Players[g.Turn].Name) {
					return errors.New("drawing failed")
				}
				curNode = curNode.Next("draw", g.Turn, len(gp.hands))
				nodesInSimulation = append(nodesInSimulation, curNode)
			} else {
				cfSelection := types.RandomInitialFoot(g)
				if err := g.Pass(ctx, g.Players[g.Turn].Name, cfSelection.X, cfSelection.Y); err != nil {
					return fmt.Errorf("passing: %w", err)
				}
				move := fmt.Sprintf("pass(%d,%d)", cfSelection.X, cfSelection.Y)
				curNode = curNode.Next(move, g.Turn, len(gp.hands))
				nodesInSimulation = append(nodesInSimulation, curNode)
			}
		} else {
			spacer := legalSpacers[whichMove-len(legalMoves)]
			if err := g.LaySpacer(ctx, g.Players[g.Turn].Name, spacer); err != nil {
				return fmt.Errorf("spacing: %w", err)
			}
			curNode = curNode.Next(spacer.String(), g.Turn, len(gp.hands))
			nodesInSimulation = append(nodesInSimulation, curNode)
		}
	}

	// The rest is fast so we still do it if we ran out of time.

	// Now that we've simulated a random game, let's backprop its eval.

	// First we start with the score at the end of this simulation.
	lastNode := nodesInSimulation[len(nodesInSimulation)-1]
	lastNode.Eval = make([]float64, len(gp.hands))
	for i := range lastNode.Eval {
		lastNode.Eval[i] = float64(g.Players[i].Score)
	}
	log.Printf("lastNode.Eval: %v", lastNode.Eval)

	// For every node, we assume that the player is maximizing their own score.
	// This is myopic to the round, missing outcomes that might win the game. Alas.
	for i := len(nodesInSimulation) - 2; i >= 0; i-- {
		cur := nodesInSimulation[i]
		bestMove, err := cur.ChooseBestMove(ctx)
		if err != nil {
			return fmt.Errorf("choosing best move: %w", err)
		}
		cur.Eval = cur.Moves[bestMove].Eval
		log.Printf("propagated eval: %v", cur.Eval)
	}
	return nil
}
