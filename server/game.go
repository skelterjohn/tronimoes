package server

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
)

type InMemoryGames struct {
	gamesMu sync.Mutex
	games   map[string]*spb.Game
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
	g.games[gm.GetGameId()] = gm
	return nil
}

func (g *InMemoryGames) ReadGame(ctx context.Context, id string) (*spb.Game, error) {
	g.gamesMu.Lock()
	defer g.gamesMu.Unlock()
	if gm, ok := g.games[id]; ok {
		return gm, nil
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
