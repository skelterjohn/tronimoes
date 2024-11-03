package game

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type Store interface {
	ReadGame(ctx context.Context, code string) (*Game, error)
	WriteGame(ctx context.Context, game *Game) error
	WatchGame(ctx context.Context, code string, version int64) <-chan *Game
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

	game, ok := s.games[code]
	if !ok {
		return nil, ErrNoSuchGame
	}

	if game.Done {
		delete(s.games, code)
		return nil, ErrNoSuchGame
	}

	// Deep copy using JSON marshal/unmarshal so that changes (like filtering)
	// aren't reflected in the saved state.
	data, err := json.Marshal(game)
	if err != nil {
		return nil, fmt.Errorf("marshaling game: %w", err)
	}

	var gameCopy Game
	if err := json.Unmarshal(data, &gameCopy); err != nil {
		return nil, fmt.Errorf("unmarshaling game: %w", err)
	}

	return &gameCopy, nil
}

func (s *MemoryStore) WriteGame(ctx context.Context, game *Game) error {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()

	game.Version++

	// Deep copy using JSON marshal/unmarshal so that changes (like filtering)
	// aren't reflected in the saved state.
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("marshaling game: %w", err)
	}

	var gameCopy Game
	if err := json.Unmarshal(data, &gameCopy); err != nil {
		return fmt.Errorf("unmarshaling game: %w", err)
	}

	s.games[game.Code] = &gameCopy

	s.watchMu.Lock()
	for _, ch := range s.watchChans[gameCopy.Code] {
		// We make yet another copy because we may get concurrent watch-reads.
		var gameCopyCopy Game
		if err := json.Unmarshal(data, &gameCopyCopy); err != nil {
			return fmt.Errorf("unmarshaling game: %w", err)
		}
		ch <- &gameCopyCopy
		close(ch)
	}
	s.watchChans[gameCopy.Code] = nil
	s.watchMu.Unlock()

	return nil
}

func (s *MemoryStore) WatchGame(ctx context.Context, code string, version int64) <-chan *Game {
	// todo race with game updated
	s.watchMu.Lock()
	defer s.watchMu.Unlock()
	ch := make(chan *Game, 1)
	s.watchChans[code] = append(s.watchChans[code], ch)
	return ch
}
