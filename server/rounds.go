package server

import (
	"context"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/tiles"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Rounds struct {
	Games Games
}

func (r *Rounds) StartRound(ctx context.Context, g *spb.Game) error {
	switch g.GetBoardShape() {
	case spb.BoardShape_standard_31_by_30:
		b := &tpb.Board{
			Width:  31,
			Height: 30,
		}
		for _, p := range g.GetPlayers() {
			b.Players = append(b.Players, &tpb.Player{
				PlayerId: p.GetPlayerId(),
			})
		}

		b, err := tiles.SetupBoard(ctx, b, 100)
		if err != nil {
			return annotatef(err, "could not set up initial board")
		}

		if err := r.Games.WriteBoard(ctx, g.GetGameId(), b); err != nil {
			return annotatef(err, "could not write board")
		}

		return nil
	}
	return status.Error(codes.InvalidArgument, "board shape not defined")
}
