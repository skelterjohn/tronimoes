package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type QueuedPlayer struct {
	Req         *tpb.CreateGameRequest
	OperationID string
}

type InMemoryQueue struct {
	JoinRequests []*QueuedPlayer
}

func (q *InMemoryQueue) AddPlayer(ctx context.Context, req *tpb.CreateGameRequest, operationID string) error {
	return status.Error(codes.NotFound, "no games ready")
}

func (q *InMemoryQueue) FindGame(ctx context.Context) (*tpb.Game, error) {
	return nil, status.Error(codes.NotFound, "no games ready")
}
