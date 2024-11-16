package game

import (
	"context"
	"encoding/json"
	"fmt"

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
		"code_prefix": game.Code[:6],
		"open":        open,
		"done":        game.Done,
		"game_json":   string(gameData),
		"version":     game.Version,
	}); err != nil {
		return fmt.Errorf("could not write: %v", err)
	}
	return nil
}
func (s *FireStore) WatchGame(ctx context.Context, code string, version int64) <-chan *Game {
	updates := make(chan *Game)

	go func() {
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

			docVersion, _ := data["version"].(int64)
			if docVersion <= version {
				continue
			}

			if gameData, ok := data["game_json"].(string); ok {
				g := &Game{}
				if err := json.Unmarshal([]byte(gameData), g); err != nil {
					continue
				}
				updates <- g
				return
			}
		}
	}()

	return updates
}

func (s *FireStore) players(ctx context.Context) *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("players")
}

func (s *FireStore) RegisterPlayerName(ctx context.Context, playerID, playerName string) error {
	if n, err := s.GetPlayerName(ctx, playerID); err == nil {
		return fmt.Errorf("already registered as %q", n)
	}

	doc, err := s.players(ctx).Doc(playerID).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("could not read player: %v", err)
	}
	if doc.Exists() {
		registeredID := doc.Data()["id"].(string)
		if registeredID != playerID {
			return fmt.Errorf("%q already registered to someone else", playerName)
		}
		return nil
	}
	_, err = s.players(ctx).Doc(playerID).Set(ctx, map[string]any{
		"name": playerName,
		"id":   playerID,
	})
	return err
}

func (s *FireStore) GetPlayerName(ctx context.Context, playerID string) (string, error) {
	iter := s.players(ctx).Where("id", "==", playerID).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return "", fmt.Errorf("could not query: %v", err)
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("no player registered")
	}
	return docs[0].Data()["name"].(string), nil
}
