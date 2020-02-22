package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

func annotatef(err error, format string, items ...interface{}) error {
	msg := fmt.Sprintf(format, items...)
	return status.Error(status.Code(err), msg)
}

type Operations interface {
	WriteOperation(ctx context.Context, op *tpb.Operation) error
	ReadOperation(ctx context.Context, id string) (*tpb.Operation, error)
	NewOperation(ctx context.Context) (*tpb.Operation, error)
}

type GameQueue interface {
	AddPlayer(ctx context.Context, req *tpb.CreateGameRequest, operationID string) error
	FindGame(ctx context.Context) (*tpb.Game, error)
}

type Tronimoes struct {
	Ops   Operations
	Queue GameQueue
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *tpb.CreateGameRequest) (*tpb.Operation, error) {
	op, err := t.Ops.NewOperation(ctx)
	if err != nil {
		return nil, annotatef(err, "could not create operation")
	}
	if err := t.Queue.AddPlayer(ctx, req, op.GetOperationId()); err != nil {
		return nil, annotatef(err, "could not create queue player")
	}
	return op, nil
}

func (t *Tronimoes) GetOperation(ctx context.Context, req *tpb.GetOperationRequest) (*tpb.Operation, error) {
	return t.Ops.ReadOperation(ctx, req.GetOperationId())
}

func (t *Tronimoes) GetGame(ctx context.Context, req *tpb.GetGameRequest) (*tpb.Game, error) {
	return &tpb.Game{
		GameId: "abc123",
	}, nil
}
