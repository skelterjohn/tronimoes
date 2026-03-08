package gibbs_planner

import (
	"context"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

func (gp *GibbsPlanner) Heuristic(ctx context.Context, g *game.Game, root *PlanNode) []float64 {
	playerLines := g.CurrentRound(ctx).PlayerLines
	h := make([]float64, len(gp.hands))
	for i := range h {
		r := float64(g.Players[i].Score) - root.Eval[i]
		// lower the score for each tile in the hand
		extraTilesPenalty := float64(len(gp.hands[i].tiles)) * -0.05
		// raise the score for each tile played
		playedTilesBonus := float64(len(playerLines[g.Players[i].Name])) * 0.05
		h[i] = r + extraTilesPenalty + playedTilesBonus
	}
	return h
}
