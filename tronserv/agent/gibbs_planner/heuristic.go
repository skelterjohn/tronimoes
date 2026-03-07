package gibbs_planner

import (
	"context"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

func (gp *GibbsPlanner) Heuristic(ctx context.Context, g *game.Game, root *PlanNode) []float64 {
	h := make([]float64, len(gp.hands))
	for i := range h {
		h[i] = float64(g.Players[i].Score) - root.Eval[i] + 1
	}
	return h
}
