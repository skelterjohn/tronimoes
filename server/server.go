package server

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
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
	MakeNextGame(ctx context.Context) error
}

type Tronimoes struct {
	Operations Operations
	GameQueue  GameQueue
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *tpb.CreateGameRequest) (*tpb.Operation, error) {
	op, err := t.Operations.NewOperation(ctx)
	if err != nil {
		return nil, annotatef(err, "could not create operation")
	}
	if err := t.GameQueue.AddPlayer(ctx, req, op.GetOperationId()); err != nil {
		return nil, annotatef(err, "could not create queue player")
	}

	if err := t.GameQueue.MakeNextGame(ctx); err != nil && status.Code(err) != codes.NotFound {
		log.Printf("Error finding the next game: %v", err)
	}

	return op, nil
}

func (t *Tronimoes) GetOperation(ctx context.Context, req *tpb.GetOperationRequest) (*tpb.Operation, error) {
	return t.Operations.ReadOperation(ctx, req.GetOperationId())
}

func (t *Tronimoes) GetGame(ctx context.Context, req *tpb.GetGameRequest) (*tpb.Game, error) {
	return &tpb.Game{
		GameId: "abc123",
	}, nil
}
