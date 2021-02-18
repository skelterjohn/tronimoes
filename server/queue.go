package server

import (
	"context"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
)

type Games interface {
	WriteGame(ctx context.Context, gm *spb.Game) error
	ReadGame(ctx context.Context, id string) (*spb.Game, error)
	NewGame(ctx context.Context, gm *spb.Game) (*spb.Game, error)
}

type queuedPlayer struct {
	PlayerID    string
	Req         *spb.CreateGameRequest
	OperationID string
}

type InMemoryQueue struct {
	mu           sync.Mutex
	joinRequests []*queuedPlayer

	Games      Games
	Operations Operations
}

func (q *InMemoryQueue) AddPlayer(ctx context.Context, playerID string, req *spb.CreateGameRequest, operationID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.joinRequests = append(q.joinRequests, &queuedPlayer{
		PlayerID:    playerID,
		Req:         req,
		OperationID: operationID,
	})
	return nil
}

func (q *InMemoryQueue) MakeNextGame(ctx context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.joinRequests) <= 1 {
		return status.Error(codes.NotFound, "no games ready")
	}

	g := &spb.Game{}
	g.PlayerIds = []string{
		q.joinRequests[0].PlayerID,
		q.joinRequests[1].PlayerID,
	}
	opIDs := []string{
		q.joinRequests[0].OperationID,
		q.joinRequests[1].OperationID,
	}

	var err error
	if g, err = q.Games.NewGame(ctx, g); err != nil {
		return annotatef(err, "could not create new game")
	}

	log.Printf("Created new game %q for %q", g.GameId, g.PlayerIds)

	ops := []*spb.Operation{}

	for _, opID := range opIDs {
		op, err := q.Operations.ReadOperation(ctx, opID)
		if err != nil {
			return annotatef(err, "could not read new player operation")
		}
		ops = append(ops, op)
	}

	gdata, err := proto.Marshal(g)
	if err != nil {
		return annotatef(err, "could not marshal game for operation payload")
	}

	for _, op := range ops {
		op.Done = true
		if err := q.Operations.WriteOperation(ctx, op); err != nil {
			log.Printf("Could not write new player operation: %v", err)
		}
		op.Payload = &any.Any{
			TypeUrl: "skelterjohn.tronimoes.Game",
			Value:   gdata,
		}
		op.Status = spb.Operation_SUCCESS
	}

	q.joinRequests = q.joinRequests[2:]

	return nil
}
