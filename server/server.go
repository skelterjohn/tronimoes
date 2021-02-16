package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
)

func annotatef(err error, format string, items ...interface{}) error {
	msg := fmt.Sprintf(format, items...)
	return status.Error(status.Code(err), msg)
}

type Operations interface {
	WriteOperation(ctx context.Context, op *spb.Operation) error
	ReadOperation(ctx context.Context, id string) (*spb.Operation, error)
	NewOperation(ctx context.Context) (*spb.Operation, error)
}

type GameQueue interface {
	AddPlayer(ctx context.Context, req *spb.CreateGameRequest, operationID string) error
	MakeNextGame(ctx context.Context) error
}

type Tronimoes struct {
	Operations Operations
	GameQueue  GameQueue
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *spb.CreateGameRequest) (*spb.Operation, error) {
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

func (t *Tronimoes) GetOperation(ctx context.Context, req *spb.GetOperationRequest) (*spb.Operation, error) {
	return t.Operations.ReadOperation(ctx, req.GetOperationId())
}

func (t *Tronimoes) GetGame(ctx context.Context, req *spb.GetGameRequest) (*spb.Game, error) {
	return &spb.Game{
		GameId: "abc123",
	}, nil
}

func Serve(ctx context.Context, port string, s *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":"+port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	operations := &InMemoryOperations{}

	tronimoes := &Tronimoes{
		Operations: operations,
		GameQueue: &InMemoryQueue{
			Games:      &InMemoryGames{},
			Operations: operations,
		},
	}

	spb.RegisterTronimoesServer(s, tronimoes)
	reflection.Register(s)

	return s.Serve(lis)
}
