package server

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type queuedPlayer struct {
	Req         *tpb.CreateGameRequest
	OperationID string
}

type InMemoryQueue struct {
	mu           sync.Mutex
	joinRequests []*queuedPlayer
}

func (q *InMemoryQueue) AddPlayer(ctx context.Context, req *tpb.CreateGameRequest, operationID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.joinRequests = append(q.joinRequests, &queuedPlayer{
		Req:         req,
		OperationID: operationID,
	})
	return nil
}

func (q *InMemoryQueue) FindGame(ctx context.Context) (*tpb.Game, error) {

	return nil, status.Error(codes.NotFound, "no games ready")
}
