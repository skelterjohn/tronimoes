package game

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/skelterjohn/tronimoes/tronserv/clog"
)

type PlayerConfig struct {
	Tileset string `json:"tileset" firestore:"tileset"`
}

type PlayerInfo struct {
	Name   string       `json:"name"`
	Id     string       `json:"id"`
	Config PlayerConfig `json:"config"`
}

type GameSummary struct {
	Code       string           `json:"code"`
	Scoreboard map[string]int64 `json:"scoreboard"`
	Updated    int64            `json:"updated"`
}

type ScoreboardSummary struct {
	Active  []GameSummary `json:"active"`
	Pickup  []GameSummary `json:"pickup"`
	Recent  []GameSummary `json:"recent"`
	Updated int64         `json:"updated"`
}

type Store interface {
	FindGameAlreadyPlaying(ctx context.Context, code, name string) (*Game, error)
	FindOpenGame(ctx context.Context, code string) (*Game, error)
	FindPickupGame(ctx context.Context) (*Game, error)
	ReadGame(ctx context.Context, code string) (*Game, error)
	WriteGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, code string) error
	WatchGame(ctx context.Context, code string, version int64) <-chan *Game
	RegisterPlayerName(ctx context.Context, playerID, playerName string) error
	GetPlayer(ctx context.Context, playerID string) (PlayerInfo, error)
	GetPlayerByName(ctx context.Context, playerName string) (PlayerInfo, error)
	RecordPlayerActive(ctx context.Context, code, playerName string, lastActive int64) error
	PlayerLastActive(ctx context.Context, code, playerName string) (int64, error)
	UpdatePlayerConfig(ctx context.Context, playerID string, config PlayerConfig) error
	ReportIssue(ctx context.Context, playerName string, game *Game, summary, whatHappened, whatShouldHappen, errorMessage string) error
	ListPickupGames(ctx context.Context, count int, updated int64) ([]GameSummary, error)
	ListActiveGames(ctx context.Context, count int, updated int64) ([]GameSummary, error)
	ListRecentGames(ctx context.Context, count int, updated int64) ([]GameSummary, error)
	ListScoreboards(ctx context.Context, count int, updated int64) (ScoreboardSummary, error)
}

type MemoryStore struct {
	games      map[string]*Game
	gamesMu    sync.Mutex
	watchChans map[string][]chan *Game
	watchMu    sync.Mutex
	players    map[string]PlayerInfo
	playersMu  sync.Mutex
	active     map[string]map[string]int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games:      make(map[string]*Game),
		watchChans: make(map[string][]chan *Game),
		players:    make(map[string]PlayerInfo),
		active:     make(map[string]map[string]int64),
	}
}

func (s *MemoryStore) FindOpenGame(ctx context.Context, code string) (*Game, error) {
	s.gamesMu.Lock()
	fullCode := ""
	for fc, g := range s.games {
		if g.Done {
			continue
		}
		if len(g.Rounds) > 0 {
			continue
		}
		if len(g.Players) == 6 {
			continue
		}
		if !strings.HasPrefix(fc, code) {
			continue
		}
		fullCode = fc
		break
	}
	s.gamesMu.Unlock()

	if fullCode == "" {
		return nil, ErrNoSuchGame
	}

	return s.ReadGame(ctx, fullCode)
}

func (s *MemoryStore) FindGameAlreadyPlaying(ctx context.Context, code, name string) (*Game, error) {
	s.gamesMu.Lock()
	fullCode := ""
	for fc, g := range s.games {
		if g.Done {
			continue
		}
		if !strings.HasPrefix(fc, code) {
			continue
		}
		amInIt := false
		for _, p := range g.Players {
			if p.Name == name {
				amInIt = true
			}
		}
		if !amInIt {
			continue
		}

		fullCode = fc
		break
	}
	s.gamesMu.Unlock()

	if fullCode == "" {
		return nil, ErrNoSuchGame
	}

	return s.ReadGame(ctx, fullCode)
}

func (s *MemoryStore) FindPickupGame(ctx context.Context) (*Game, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	for _, g := range s.games {
		if g.Done {
			continue
		}
		if len(g.Rounds) > 0 {
			continue
		}
		if len(g.Players) == 6 {
			continue
		}
		if !g.Pickup {
			continue
		}
		return g, nil
	}
	return nil, ErrNoSuchGame
}

func (s *MemoryStore) ReadGame(ctx context.Context, code string) (*Game, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()

	game, ok := s.games[code]
	if !ok {
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

func (s *MemoryStore) DeleteGame(ctx context.Context, code string) error {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	delete(s.games, code)
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

func (s *MemoryStore) RegisterPlayerName(ctx context.Context, playerID, playerName string) error {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()
	s.players[playerID] = PlayerInfo{
		Name: playerName,
		Id:   playerID,
	}
	return nil
}

func (s *MemoryStore) GetPlayer(ctx context.Context, playerID string) (PlayerInfo, error) {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()
	pi, ok := s.players[playerID]
	if ok {
		return pi, nil
	}
	return PlayerInfo{}, ErrNoRegisteredPlayer
}

func (s *MemoryStore) GetPlayerByName(ctx context.Context, playerName string) (PlayerInfo, error) {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()
	for _, pi := range s.players {
		if pi.Name == playerName {
			return pi, nil
		}
	}
	return PlayerInfo{}, ErrNoRegisteredPlayer
}

func (s *MemoryStore) RecordPlayerActive(ctx context.Context, code, playerName string, lastActive int64) error {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	if _, ok := s.active[code]; !ok {
		s.active[code] = make(map[string]int64)
	}
	s.active[code][playerName] = lastActive
	return nil
}
func (s *MemoryStore) PlayerLastActive(ctx context.Context, code, playerName string) (int64, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	ga, ok := s.active[code]
	if !ok {
		return 0, ErrNoSuchGame
	}
	pa, ok := ga[playerName]
	if !ok {
		return 0, ErrNoSuchPlayer
	}
	return pa, nil
}

func (s *MemoryStore) UpdatePlayerConfig(ctx context.Context, playerID string, config PlayerConfig) error {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()
	p := s.players[playerID]
	p.Config = config
	s.players[playerID] = p
	return nil
}

var filenameSafeRE = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func (s *MemoryStore) ReportIssue(ctx context.Context, playerName string, game *Game, summary, whatHappened, whatShouldHappen, errorMessage string) error {
	clog.Info(ctx, "reporting issue for game", "code", game.Code)
	data, err := json.MarshalIndent(game, "", "\t")
	if err != nil {
		return fmt.Errorf("marshaling game: %w", err)
	}
	name := strings.ToLower(strings.TrimSpace(summary))
	if name == "" {
		name = "issue"
	}
	name = filenameSafeRE.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		name = "issue"
	}
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable: %w", err)
	}
	testdataDir := filepath.Join(filepath.Dir(exe), "game", "testdata")
	path := filepath.Join(testdataDir, name+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	clog.Info(ctx, "wrote issue", "path", path)
	return nil
}

func SummarizeGame(g *Game) GameSummary {
	gs := GameSummary{
		Code:       g.Code,
		Scoreboard: make(map[string]int64),
	}
	for _, p := range g.Players {
		gs.Scoreboard[p.Name] = int64(p.Score)
	}
	return gs
}

func (s *MemoryStore) ListPickupGames(ctx context.Context, count int, updated int64) ([]GameSummary, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	games := make([]GameSummary, 0, count)
	for _, g := range s.games {
		started := len(g.Rounds) > 0
		if !g.Pickup || started {
			continue
		}
		games = append(games, SummarizeGame(g))
		count--
		if count <= 0 {
			break
		}
	}
	return games, nil
}

func (s *MemoryStore) ListActiveGames(ctx context.Context, count int, updated int64) ([]GameSummary, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	games := make([]GameSummary, 0, count)
	for _, g := range s.games {
		started := len(g.Rounds) > 0
		if !started {
			continue
		}
		games = append(games, SummarizeGame(g))
		count--
		if count <= 0 {
			break
		}
	}
	return games, nil
}

func (s *MemoryStore) ListRecentGames(ctx context.Context, count int, updated int64) ([]GameSummary, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	games := make([]GameSummary, 0, count)
	for _, g := range s.games {
		if !g.Done {
			continue
		}
		games = append(games, SummarizeGame(g))
		count--
		if count <= 0 {
			break
		}
	}
	return games, nil
}
func (s *MemoryStore) ListScoreboards(ctx context.Context, count int, updated int64) (ScoreboardSummary, error) {
	active, err := s.ListActiveGames(ctx, count, updated)
	if err != nil {
		return ScoreboardSummary{}, err
	}
	pickup, err := s.ListPickupGames(ctx, count, updated)
	if err != nil {
		return ScoreboardSummary{}, err
	}
	recent, err := s.ListRecentGames(ctx, count, updated)
	if err != nil {
		return ScoreboardSummary{}, err
	}

	return ScoreboardSummary{
		Active:  active,
		Pickup:  pickup,
		Recent:  recent,
		Updated: updated,
	}, nil
}
