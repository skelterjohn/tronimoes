package types

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type Move struct {
	LayTile     bool
	LaidTile    game.LaidTile
	PlaceSpacer bool
	Spacer      game.Spacer
	Draw        bool
	Pass        bool
	Selected    game.Coord
}

func (m Move) String() string {
	if m.LayTile {
		return fmt.Sprintf("LT{%s %s %s %s %d %v}", m.LaidTile, m.LaidTile.PlayerName, m.LaidTile.WhoLaidIt, m.LaidTile.Orientation, m.LaidTile.NextPips, m.LaidTile.Dead)
	}
	if m.PlaceSpacer {
		return fmt.Sprintf("PS{%s}", m.Spacer)
	}
	if m.Draw {
		return "Draw"
	}
	if m.Pass {
		return fmt.Sprintf("Pass{%v}", m.Selected)
	}
	return ""
}

func InferMove(ctx context.Context, pg *game.Game, g *game.Game) (Move, bool) {
	if len(g.Rounds) == 0 {
		return Move{}, false
	}
	currentRound := g.Rounds[len(g.Rounds)-1]
	previousCurrentRound := pg.CurrentRound(ctx)
	if currentRound == nil || previousCurrentRound == nil {
		return Move{}, false
	}
	lastPlayer := g.Players[pg.Turn]
	if len(currentRound.LaidTiles) > len(previousCurrentRound.LaidTiles) {
		return Move{
			LayTile:  true,
			LaidTile: *currentRound.LaidTiles[len(currentRound.LaidTiles)-1],
		}, true
	}
	if currentRound.Spacer != nil {
		return Move{
			PlaceSpacer: true,
			Spacer:      *currentRound.Spacer,
		}, true
	}
	for _, p := range g.Players {
		if p.Name != lastPlayer.Name {
			continue
		}
		if len(p.Hand) > len(pg.GetPlayer(ctx, p.Name).Hand) {
			return Move{
				Draw: true,
			}, true
		}
	}
	if pg.Turn != g.Turn {
		selected := game.Coord{X: -1, Y: -1}
		if len(currentRound.PlayerLines[lastPlayer.Name]) == 1 {
			selected = lastPlayer.ChickenFootCoord
		}
		return Move{
			Pass:     true,
			Selected: selected,
		}, true
	}
	return Move{}, false
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
