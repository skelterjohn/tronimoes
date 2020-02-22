package server

import (
	"context"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Tronimoes struct {
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *tpb.CreateGameRequest) (op *tpb.Operation, err error) {
	return &tpb.Operation{
		OperationId: "abc123",
	}, nil
}

func (t *Tronimoes) GetOperation(ctx context.Context, req *tpb.GetOperationRequest) (op *tpb.Operation, err error) {
	return &tpb.Operation{
		OperationId: "abc123",
	}, nil
}

func (t *Tronimoes) GetGame(ctx context.Context, req *tpb.GetGameRequest) (g *tpb.Game, err error) {
	return &tpb.Game{
		GameId: "abc123",
	}, nil
}
