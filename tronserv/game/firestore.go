package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FireStore struct {
	storeClient *firestore.Client
	env         string
}

func NewFirestore(ctx context.Context, project, env string) (*FireStore, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: project,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing firestore app: %v", err)
	}
	storeClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing firestore storeClient: %v", err)
	}
	return &FireStore{
		storeClient: storeClient,
		env:         env,
	}, nil
}

func (s *FireStore) games(ctx context.Context) *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("games")
}

func (s *FireStore) FindGameAlreadyPlaying(ctx context.Context, code, name string) (*Game, error) {
	c := s.games(ctx)
	iter := c.Where("code_prefix", "==", code).Where("done", "==", false).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}

	if len(docs) == 0 {
		return nil, ErrNoSuchGame
	}

	for _, doc := range docs {
		data := doc.Data()
		gameData, ok := data["game_json"].(string)
		if !ok {
			return nil, fmt.Errorf("bad data type for game_json: %T", data["game_json"])
		}

		g := &Game{}
		if err := json.Unmarshal([]byte(gameData), g); err != nil {
			return nil, fmt.Errorf("could not unmarshal: %v", err)
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

		return g, nil
	}
	return nil, nil
}

func (s *FireStore) FindOpenGame(ctx context.Context, code string) (*Game, error) {
	c := s.games(ctx)
	iter := c.Where("code_prefix", "==", code).Where("open", "==", true).Where("done", "==", false).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}

	if len(docs) == 0 {
		return nil, ErrNoSuchGame
	}

	// Return the first matching game
	doc := docs[0]
	data := doc.Data()
	gameData, ok := data["game_json"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type for game_json: %T", data["game_json"])
	}

	g := &Game{}
	if err := json.Unmarshal([]byte(gameData), g); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %v", err)
	}

	return g, nil
}

func (s *FireStore) FindPickupGame(ctx context.Context) (*Game, error) {
	c := s.games(ctx)
	iter := c.Where("open", "==", true).Where("done", "==", false).Where("pickup", "==", true).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	if len(docs) == 0 {
		return nil, ErrNoSuchGame
	}
	doc := docs[0]
	data := doc.Data()
	gameData, ok := data["game_json"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type for game_json: %T", data["game_json"])
	}

	g := &Game{}
	if err := json.Unmarshal([]byte(gameData), g); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %v", err)
	}

	return g, nil
}

func (s *FireStore) ReadGame(ctx context.Context, code string) (*Game, error) {
	c := s.games(ctx)
	doc, err := c.Doc(code).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNoSuchGame
		}
		return nil, fmt.Errorf("could not read: %v", err)
	}

	data := doc.Data()
	gameData, ok := data["game_json"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type for game_json: %T", data["game_json"])
	}

	g := &Game{}
	if err := json.Unmarshal([]byte(gameData), g); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %v", err)
	}

	return g, nil
}
func (s *FireStore) WriteGame(ctx context.Context, game *Game) error {
	open := len(game.Rounds) == 0 && len(game.Players) < 6

	game.Version++
	gameData, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("could not marshal: %v", err)
	}
	c := s.games(ctx)
	if _, err := c.Doc(game.Code).Set(ctx, map[string]any{
		"created":     game.Created,
		"code_prefix": game.Code[:6],
		"open":        open,
		"pickup":      game.Pickup,
		"done":        game.Done,
		"game_json":   string(gameData),
		"version":     game.Version,
	}); err != nil {
		return fmt.Errorf("could not write: %v", err)
	}
	return nil
}

func (s *FireStore) DeleteGame(ctx context.Context, code string) error {
	_, err := s.games(ctx).Doc(code).Delete(ctx)
	return err
}

func (s *FireStore) WatchGame(ctx context.Context, code string, version int64) <-chan *Game {
	updates := make(chan *Game)

	go func(ctx context.Context) {
		defer close(updates)

		iter := s.games(ctx).Doc(code).Snapshots(ctx)
		for {
			snap, err := iter.Next()
			if err != nil {
				return
			}

			data := snap.Data()
			if data == nil {
				continue
			}

			docVersion, ok := data["version"].(int64)
			if !ok {
				log.Printf("bad data type for version: %T", data["version"])
				continue
			}
			if docVersion <= version {
				continue
			}

			if gameData, ok := data["game_json"].(string); ok {
				g := &Game{}
				if err := json.Unmarshal([]byte(gameData), g); err != nil {
					log.Printf("could not unmarshal: %v", err)
					continue
				}
				updates <- g
				return
			}
		}
	}(ctx)

	return updates
}

func (s *FireStore) players(ctx context.Context) *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("players")
}

func (s *FireStore) RegisterPlayerName(ctx context.Context, playerID, playerName string) error {
	if pi, err := s.GetPlayer(ctx, playerID); err == nil {
		return fmt.Errorf("already registered as %q", pi.Name)
	}
	_, err := s.players(ctx).Doc(playerID).Set(ctx, map[string]any{
		"name": playerName,
		"id":   playerID,
	})
	return err
}

func (s *FireStore) GetPlayer(ctx context.Context, playerID string) (PlayerInfo, error) {
	doc, err := s.players(ctx).Doc(playerID).Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return PlayerInfo{}, ErrNoRegisteredPlayer
	}
	if err != nil {
		return PlayerInfo{}, fmt.Errorf("could not read: %v", err)
	}
	pi := PlayerInfo{}
	if name, ok := doc.Data()["name"].(string); ok {
		pi.Name = name
	}
	if id, ok := doc.Data()["id"].(string); ok {
		pi.Id = id
	}

	return pi, nil
}

func (s *FireStore) GetPlayerByName(ctx context.Context, playerName string) (PlayerInfo, error) {
	iter := s.players(ctx).Where("name", "==", playerName).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return PlayerInfo{}, fmt.Errorf("could not query: %v", err)
	}
	if len(docs) == 0 {
		return PlayerInfo{}, ErrNoRegisteredPlayer
	}
	pi := PlayerInfo{}
	if name, ok := docs[0].Data()["name"].(string); ok {
		pi.Name = name
	}
	if id, ok := docs[0].Data()["id"].(string); ok {
		pi.Id = id
	}

	return pi, nil
}

func (s *FireStore) RecordPlayerActive(ctx context.Context, code, playerName string, lastActive int64) error {
	_, err := s.games(ctx).Doc(code).Collection("active").Doc(playerName).Set(ctx, map[string]any{
		"last_active": lastActive,
	})
	return err
}

func (s *FireStore) PlayerLastActive(ctx context.Context, code, playerName string) (int64, error) {
	doc, err := s.games(ctx).Doc(code).Collection("active").Doc(playerName).Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not read: %v", err)
	}
	if lastActive, ok := doc.Data()["last_active"].(int64); ok {
		return lastActive, nil
	}
	return 0, fmt.Errorf("bad data type for last_active: %T", doc.Data()["last_active"])
}

func (s *FireStore) UpdatePlayerConfig(ctx context.Context, playerID string, config PlayerConfig) error {
	_, err := s.players(ctx).Doc(playerID).Set(ctx, map[string]any{"config": config})
	return err
}
