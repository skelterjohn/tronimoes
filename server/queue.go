package server

import (
	"context"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Games interface {
	WriteGame(ctx context.Context, gm *tpb.Game) error
	ReadGame(ctx context.Context, id string) (*tpb.Game, error)
	NewGame(ctx context.Context, gm *tpb.Game) (*tpb.Game, error)
}

type queuedPlayer struct {
	Req         *tpb.CreateGameRequest
	OperationID string
}

type InMemoryQueue struct {
	mu           sync.Mutex
	joinRequests []*queuedPlayer

	Games      Games
	Operations Operations
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

func (q *InMemoryQueue) MakeNextGame(ctx context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.joinRequests) <= 1 {
		return status.Error(codes.NotFound, "no games ready")
	}

	g := &tpb.Game{}
	g.Players = []string{
		q.joinRequests[0].Req.GetPlayerSelf(),
		q.joinRequests[1].Req.GetPlayerSelf(),
	}
	opIDs := []string{
		q.joinRequests[0].OperationID,
		q.joinRequests[1].OperationID,
	}

	var err error
	if g, err = q.Games.NewGame(ctx, g); err != nil {
		return annotatef(err, "could not create new game")
	}

	log.Printf("Created new game %q for %q", g.GameId, g.Players)

	ops := []*tpb.Operation{}

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
	}

	q.joinRequests = q.joinRequests[2:]

	return nil
}
