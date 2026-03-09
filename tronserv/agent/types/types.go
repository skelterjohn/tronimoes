package types

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type Move struct {
	LaidTile *game.LaidTile
	Spacer   *game.Spacer
	Draw     bool
	Pass     bool
	Selected game.Coord
}

func (m Move) String() string {
	return fmt.Sprintf("Move{%v, %v, %v, %v, %v}", m.LaidTile, m.Spacer, m.Draw, m.Pass, m.Selected)
}

type Agent interface {
	Ready(ctx context.Context)
	Update(ctx context.Context, previousGame *game.Game, g *game.Game)
	GetMove(ctx context.Context, g *game.Game, p *game.Player) Move
	CompleteRound(ctx context.Context, g *game.Game)
	CompleteGame(ctx context.Context, g *game.Game)
}

func RandomInitialFoot(g *game.Game) game.Coord {
	cfSelection := game.Coord{
		X: g.BoardWidth / 2,
		Y: (g.BoardHeight / 2),
	}
	var dx = rand.Intn(2)
	dy := rand.Intn(3) - 1
	if dy == 0 {
		if dx == 0 {
			dx = -1
		} else {
			dx = 2
		}
	}
	cfSelection.X += dx
	cfSelection.Y += dy
	return cfSelection
}
