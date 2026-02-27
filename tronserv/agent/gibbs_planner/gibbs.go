package gibbs_planner

import (
	"context"
	"fmt"
	"math/rand"

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
