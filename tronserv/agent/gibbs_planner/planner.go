package gibbs_planner

import (
	"context"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type PlanNode struct {
	Turn  int
	Moves map[string]PlanNode
	Eval  float64
}

func (gp *GibbsPlanner) Plan(ctx context.Context, g *game.Game) {
	root := &PlanNode{
		Turn:  g.Turn,
		Moves: make(map[string]PlanNode),
	}
	_ = root
}
