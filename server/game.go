package server

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

type InMemoryGames struct {
	gamesMu  sync.Mutex
	games    map[string]*spb.Game
	boardsMu sync.Mutex
	boards   map[string]*tpb.Board
}

func (g *InMemoryGames) WriteGame(ctx context.Context, gm *spb.Game) error {
	if gm.GetGameId() == "" {
		return status.Error(codes.InvalidArgument, "no game ID")
	}
	g.gamesMu.Lock()
	defer g.gamesMu.Unlock()
	if g.games == nil {
		g.games = map[string]*spb.Game{}
	}
	g.games[gm.GetGameId()] = proto.Clone(gm).(*spb.Game)
	return nil
}

func (g *InMemoryGames) ReadGame(ctx context.Context, id string) (*spb.Game, error) {
	g.gamesMu.Lock()
	defer g.gamesMu.Unlock()
	if gm, ok := g.games[id]; ok {
		return proto.Clone(gm).(*spb.Game), nil
	}
	return nil, status.Errorf(codes.NotFound, "no such game %s", id)
}

func (g *InMemoryGames) NewGame(ctx context.Context, gm *spb.Game) (*spb.Game, error) {
	gm.GameId = uuid.New().String()
	if err := g.WriteGame(ctx, gm); err != nil {
		return nil, err
	}
	return gm, nil
}

func (g *InMemoryGames) WriteBoard(ctx context.Context, id string, b *tpb.Board) error {
	g.boardsMu.Lock()
	defer g.boardsMu.Unlock()
	if g.boards == nil {
		g.boards = make(map[string]*tpb.Board)
	}
	g.boards[id] = proto.Clone(b).(*tpb.Board)
	return nil
}

func (g *InMemoryGames) ReadBoard(ctx context.Context, id string) (*tpb.Board, error) {
	g.boardsMu.Lock()
	defer g.boardsMu.Unlock()
	if b, ok := g.boards[id]; ok {
		return proto.Clone(b).(*tpb.Board), nil
	}
	return nil, status.Errorf(codes.NotFound, "no board for game %s", id)
}
