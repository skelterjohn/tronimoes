package types

import (
	"context"
	"fmt"

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
}
