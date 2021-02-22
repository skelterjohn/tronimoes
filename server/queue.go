package server

import (
	"context"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/tiles"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

type Games interface {
	WriteGame(ctx context.Context, gm *spb.Game) error
	ReadGame(ctx context.Context, id string) (*spb.Game, error)
	WriteBoard(ctx context.Context, id string, b *tpb.Board) error
	ReadBoard(ctx context.Context, id string) (*tpb.Board, error)
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

	jrs := []*queuedPlayer{
		q.joinRequests[0],
		q.joinRequests[1],
	}

	g := &spb.Game{}
	opIDs := []string{}
	for _, jr := range jrs {
		g.PlayerIds = append(g.PlayerIds, jr.PlayerID)
		opIDs = append(opIDs, jr.OperationID)
	}

	g.GameId = uuid.New().String()

	if err := q.Games.WriteGame(ctx, g); err != nil {
		return annotatef(err, "could not write new game")
	}

	switch q.joinRequests[0].Req.GetBoardShape() {
	case spb.CreateGameRequest_standard_31_by_30:
		b := &tpb.Board{
			Width:  31,
			Height: 30,
		}
		for _, pid := range g.PlayerIds {
			b.Players = append(b.Players, &tpb.Player{
				PlayerId: pid,
			})
		}

		b, err := tiles.SetupBoard(ctx, b, 100)
		if err != nil {
			return annotatef(err, "could not set up initial board")
		}

		if err := q.Games.WriteBoard(ctx, g.GetGameId(), b); err != nil {
			return annotatef(err, "could not write board")
		}
	default:
		return status.Error(codes.FailedPrecondition, "board shape not defined")
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
