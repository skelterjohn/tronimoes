package server

import (
	"context"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Tronimoes struct {
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *tpb.CreateGameRequest) (*tpb.Operation, error) {
	return &tpb.Operation{
		OperationId: "abc123",
	}, nil
}

func (t *Tronimoes) GetOperation(ctx context.Context, req *tpb.GetOperationRequest) (*tpb.Operation, error) {
	return &tpb.Operation{
		OperationId: "abc123",
	}, nil
}

func (t *Tronimoes) GetGame(ctx context.Context, req *tpb.GetGameRequest) (*tpb.Game, error) {
	return &tpb.Game{
		GameId: "abc123",
	}, nil
}
