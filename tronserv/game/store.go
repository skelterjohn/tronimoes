package game

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoSuchGame = errors.New("no such game")
)

type Store interface {
	ReadGame(ctx context.Context, code string) (*Game, error)
	WriteGame(ctx context.Context, game *Game) error
	WatchGame(ctx context.Context, code string) <-chan *Game
}

type MemoryStore struct {
	games      map[string]*Game
	gamesMu    sync.Mutex
	watchChans map[string][]chan *Game
	watchMu    sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games:      make(map[string]*Game),
		watchChans: make(map[string][]chan *Game),
	}
}

func (s *MemoryStore) ReadGame(ctx context.Context, code string) (*Game, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	if game, ok := s.games[code]; ok {
		return game, nil
	}
	return nil, nil
}

func (s *MemoryStore) WriteGame(ctx context.Context, game *Game) error {
	game.Version++

	s.gamesMu.Lock()
	s.games[game.Code] = game
	s.gamesMu.Unlock()

	s.watchMu.Lock()
	for _, ch := range s.watchChans[game.Code] {
		ch <- game
		close(ch)
	}
	s.watchChans[game.Code] = nil
	s.watchMu.Unlock()

	return nil
}

func (s *MemoryStore) WatchGame(ctx context.Context, code string) <-chan *Game {
	s.watchMu.Lock()
	defer s.watchMu.Unlock()
	ch := make(chan *Game, 1)
	s.watchChans[code] = append(s.watchChans[code], ch)
	return ch
}
