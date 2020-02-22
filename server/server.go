package server

import (
	"context"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Operations interface {
	WriteOperation(ctx context.Context, op *tpb.Operation) error
	ReadOperation(ctx context.Context, id string) (*tpb.Operation, error)
	NewOperation(ctx context.Context) (*tpb.Operation, error)
}

type Tronimoes struct {
	Ops Operations
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *tpb.CreateGameRequest) (op *tpb.Operation, err error) {
	return t.Ops.NewOperation(ctx)
}

func (t *Tronimoes) GetOperation(ctx context.Context, req *tpb.GetOperationRequest) (op *tpb.Operation, err error) {
	return t.Ops.ReadOperation(ctx, req.OperationId)
}

func (t *Tronimoes) GetGame(ctx context.Context, req *tpb.GetGameRequest) (g *tpb.Game, err error) {
	return &tpb.Game{
		GameId: "abc123",
	}, nil
}
