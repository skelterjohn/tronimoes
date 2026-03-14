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
	Moves   map[types.Move]*PlanNode
	Eval    []float64
	R       []float64
	V       []float64
	H       []float64
}

func NewPlanNode(turn int, count int, depth int) *PlanNode {
	return &PlanNode{
		Turn:  turn,
		Depth: depth,
		Moves: make(map[types.Move]*PlanNode),
		Eval:  make([]float64, count),
		R:     make([]float64, count),
		V:     make([]float64, count),
		H:     make([]float64, count),
	}
}

var CFOffsets = []game.Coord{
	{X: -1, Y: 0},
	{X: 2, Y: 0},
	{X: 0, Y: 1},
	{X: 0, Y: -1},
	{X: 1, Y: 1},
	{X: 1, Y: -1},
}

func (n *PlanNode) Next(ctx context.Context, move types.Move, turn, count int) *PlanNode {
	nextNode, ok := n.Moves[move]
	if !ok {
		isDouble := false
		if move.LayTile {
			isDouble = move.LaidTile.Tile.PipsA == move.LaidTile.Tile.PipsB
		}
		if move.Pass || (move.LayTile && !isDouble) {
			turn = (turn + 1) % count
		}
		nextNode = NewPlanNode(turn, count, n.Depth+1)
		n.Moves[move] = nextNode
	}
	return nextNode
}

func (n *PlanNode) ChooseBestMove(ctx context.Context) types.Move {
	var bestMove types.Move
	bestV := -math.MaxFloat64
	for move, next := range n.Moves {
		nextV := next.V[n.Turn]
		if nextV > bestV {
			bestV = nextV
			bestMove = move
		}
	}
	return bestMove
}

func (n *PlanNode) Cull(ctx context.Context, moves map[types.Move]bool) {
	delCount := 0
	for m := range n.Moves {
		if !moves[m] {
			delete(n.Moves, m)
			delCount++
		}
	}
	if delCount > 0 {
		game.Debug(ctx, "culled %d moves, %d remain", delCount, len(n.Moves))
	}
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

		game.Debug(ctx, "curNode.depth=%d", curNode.Depth)
		cachedTileCount := 0
		cachedSpacerCount := 0
		for m := range curNode.Moves {
			if m.LayTile {
				cachedTileCount++
			} else if m.PlaceSpacer {
				cachedSpacerCount++
			}
		}
		game.Debug(ctx, "cachedTileCount=%d, cachedSpacerCount=%d", cachedTileCount, cachedSpacerCount)

		p := g.Players[g.Turn]
		playingOffRoundLeader := len(r.PlayerLines[p.Name]) == 1

		legalMoves, legalSpacers := r.FindLegalMoves(ctx, g, p)

		allMoves := map[types.Move]bool{}
		for _, lt := range legalMoves {
			allMoves[types.Move{LayTile: true, LaidTile: lt}] = true
		}
		for _, spacer := range legalSpacers {
			allMoves[types.Move{PlaceSpacer: true, Spacer: spacer}] = true
		}
		if !p.JustDrew && len(g.Bag) > 0 {
			allMoves[types.Move{Draw: true}] = true
		} else {
			if playingOffRoundLeader {
				for _, cf := range CFOffsets {
					allMoves[types.Move{Pass: true, Selected: game.Coord{X: g.BoardWidth/2 + cf.X, Y: g.BoardHeight/2 + cf.Y}}] = true
				}
			} else {
				allMoves[types.Move{Pass: true, Selected: game.Coord{X: -1, Y: -1}}] = true
			}
		}

		curNode.Cull(ctx, allMoves)

		game.Debug(ctx, "%s has %d tiles, %d spacers", p.Name, len(legalMoves), len(legalSpacers))
		moveCount := len(legalMoves) + len(legalSpacers)

		passOrDrawOptions := 1
		if playingOffRoundLeader && p.JustDrew {
			passOrDrawOptions = 6 // (all 6 ways to pass)
		}
		moveCount += passOrDrawOptions

		unnormalizedLogLikelihoods := make([]float64, moveCount)
		for i := range unnormalizedLogLikelihoods {
			unnormalizedLogLikelihoods[i] = 0
		}

		for i, lt := range legalMoves {
			nn := curNode.Next(ctx, types.Move{LayTile: true, LaidTile: lt}, g.Turn, len(gp.hands))
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

		tileMoves := len(legalMoves)
		tileAndSpacerMoves := tileMoves + len(legalSpacers)

		if whichMove < tileMoves {
			move := legalMoves[whichMove]
			game.Debug(ctx, "p%d lays %s", g.Turn, move)
			if err := g.LayTile(ctx, p.Name, &move); err != nil {
				game.Debug(ctx, "player=%+v", p)
				return fmt.Errorf("laying: %w", err)
			}
			move.Dead = false // this screws up the hashmap since FindLegalMoves doesn't set this.
			bestMove = types.Move{LayTile: true, LaidTile: move}
		} else if whichMove < tileAndSpacerMoves {
			spacer := legalSpacers[whichMove-tileMoves]
			if err := g.LaySpacer(ctx, p.Name, &spacer); err != nil {
				return fmt.Errorf("spacing: %w", err)
			}
			bestMove = types.Move{PlaceSpacer: true, Spacer: spacer}
		} else {
			if !p.JustDrew && len(g.Bag) > 0 {
				if !g.DrawTile(ctx, p.Name) {
					return errors.New("drawing failed")
				}
				bestMove = types.Move{Draw: true}
			} else {
				whichOption := whichMove - tileAndSpacerMoves
				cfSelection := game.Coord{X: -1, Y: -1}
				if playingOffRoundLeader {
					cfSelection = game.Coord{
						X: g.BoardWidth/2 + CFOffsets[whichOption].X,
						Y: g.BoardHeight/2 + CFOffsets[whichOption].Y,
					}
				}
				if err := g.Pass(ctx, p.Name, cfSelection.X, cfSelection.Y); err != nil {
					return fmt.Errorf("passing: %w", err)
				}
				bestMove = types.Move{Pass: true, Selected: cfSelection}
			}
		}

		game.Debug(ctx, "p%d -> %s", curNode.Turn, bestMove)

		nextNode = curNode.Next(ctx, bestMove, g.Turn, len(gp.hands))
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
		bestMove := cur.ChooseBestMove(ctx)
		bestNode := cur.Moves[bestMove]
		if bestNode == nil {
			// This is the last node in the simulation, so we just copy the reward.
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
