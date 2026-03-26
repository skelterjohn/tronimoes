package game

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/skelterjohn/tronimoes/tronserv/clog"
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

func (s *FireStore) games() *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("games")
}

func (s *FireStore) players() *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("players")
}

func (s *FireStore) issues() *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("issues")
}

func (s *FireStore) scoreboards() *firestore.CollectionRef {
	return s.storeClient.Collection("envs").Doc(s.env).Collection("scoreboards")
}

func (s *FireStore) FindGameAlreadyPlaying(ctx context.Context, code, name string) (*Game, error) {
	c := s.games()
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

func (s *FireStore) FindActiveCodeFromPrefix(ctx context.Context, prefix string) (string, error) {
	c := s.games()
	iter := c.Where("code_prefix", "==", prefix).Where("done", "==", false).Limit(1).Select("code_prefix").Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return "", fmt.Errorf("could not get next: %w", err)
	}
	defer iter.Stop()
	return doc.Ref.ID, nil
}

func (s *FireStore) FindOpenGame(ctx context.Context, code string) (*Game, error) {
	c := s.games()
	iter := c.Where("code_prefix", "==", code).Where("open", "==", true).Where("done", "==", false).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	if len(docs) == 0 {
		return nil, ErrNoSuchGame
	}

	// Randomly choose an open game.
	doc := docs[rand.Intn(len(docs))]
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
	c := s.games()
	iter := c.Where("open", "==", true).Where("done", "==", false).Where("pickup", "==", true).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	if len(docs) == 0 {
		return nil, ErrNoSuchGame
	}

	// Randomly choose an open game.
	doc := docs[rand.Intn(len(docs))]
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
	c := s.games()
	doc, err := c.Doc(code).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Clear the scoreboard if it exists.
			if _, err := s.scoreboards().Doc(code).Delete(ctx); err != nil {
				clog.Error(ctx, "could not delete scoreboard", err)
			}
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
	noRounds := len(game.Rounds) == 0
	open := noRounds && len(game.Players) < 6
	expectedVersion := game.Version
	game.Version++
	gameData, err := json.Marshal(game)
	if err != nil {
		game.Version--
		return fmt.Errorf("could not marshal: %v", err)
	}
	scoreboard := make(map[string]int64)
	for _, p := range game.Players {
		scoreboard[p.Name] = int64(p.Score)
	}

	c := s.games()
	docRef := c.Doc(game.Code)
	scoreboardDocRef := s.scoreboards().Doc(game.Code)
	payload := map[string]any{
		"created":     game.Created,
		"code_prefix": game.Code[:6],
		"open":        open,
		"pickup":      game.Pickup,
		"done":        game.Done,
		"game_json":   string(gameData),
		"version":     game.Version,
	}
	err = s.storeClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(docRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if doc != nil && doc.Exists() {
			data := doc.Data()
			storedVersion, ok := data["version"].(int64)
			if !ok {
				clog.Info(ctx, "unexpected version type", "type", fmt.Sprintf("%T", data["version"]))
			}
			if ok && storedVersion != expectedVersion {
				return ErrVersionConflict
			}
		}

		existingScoreboardDoc, err := tx.Get(scoreboardDocRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		shouldUpdateScoreboard := true
		if existingScoreboardDoc != nil && existingScoreboardDoc.Exists() {
			existingData := existingScoreboardDoc.Data()
			existingDone, doneOK := existingData["done"].(bool)
			existingOpen, openOK := existingData["open"].(bool)
			existingScoreboard, scoreboardOK := existingData["scoreboard"].(map[string]any)
			existingHasNonZero, hasNonZeroOK := existingData["has_nonzero_score"].(bool)
			newHasNonZero := hasNonZeroScoreboard(scoreboard)
			shouldUpdateScoreboard = !doneOK || !openOK || !scoreboardOK || !hasNonZeroOK ||
				existingDone != game.Done ||
				existingOpen != open ||
				existingHasNonZero != newHasNonZero ||
				!scoreboardMatches(scoreboard, existingScoreboard)
		}
		if err := tx.Set(docRef, payload); err != nil {
			return err
		}

		if !shouldUpdateScoreboard {
			return nil
		}

		scoreboardPayload := map[string]any{
			"scoreboard":        scoreboard,
			"has_nonzero_score": hasNonZeroScoreboard(scoreboard),
			"updated":           time.Now().Unix(),
			"open":              open,
			"pickup":            game.Pickup,
			"done":              game.Done,
		}
		return tx.Set(scoreboardDocRef, scoreboardPayload)
	})
	if err != nil {
		game.Version--
		if err == ErrVersionConflict {
			return err
		}
		return fmt.Errorf("could not write: %v", err)
	}
	return nil
}

func scoreboardMatches(next map[string]int64, current map[string]any) bool {
	if len(next) != len(current) {
		return false
	}
	for name, score := range next {
		raw, ok := current[name]
		if !ok {
			return false
		}
		existingScore, ok := raw.(int64)
		if !ok || existingScore != score {
			return false
		}
	}
	return true
}

func (s *FireStore) DeleteGame(ctx context.Context, code string) error {
	_, err := s.games().Doc(code).Delete(ctx)
	return err
}

func (s *FireStore) WatchGame(ctx context.Context, code string, version int64) <-chan *Game {
	updates := make(chan *Game)

	go func(ctx context.Context) {
		defer close(updates)

		iter := s.games().Doc(code).Snapshots(ctx)
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
				clog.Info(ctx, "bad data type for version", "type", fmt.Sprintf("%T", data["version"]))
				continue
			}
			if docVersion <= version {
				continue
			}

			if gameData, ok := data["game_json"].(string); ok {
				g := &Game{}
				if err := json.Unmarshal([]byte(gameData), g); err != nil {
					clog.Error(ctx, "could not unmarshal game", err)
					continue
				}
				updates <- g
				return
			}
		}
	}(ctx)

	return updates
}
func (s *FireStore) RegisterPlayerName(ctx context.Context, playerID, playerName string) error {
	if pi, err := s.GetPlayer(ctx, playerID); err == nil {
		return fmt.Errorf("already registered as %q", pi.Name)
	}
	_, err := s.players().Doc(playerID).Set(ctx, map[string]any{
		"name": playerName,
		"id":   playerID,
	})
	return err
}

func playerConfigFromData(data map[string]any) (PlayerConfig, error) {
	if data["config"] == nil {
		return PlayerConfig{}, nil
	}
	m, ok := data["config"].(map[string]any)
	if !ok {
		return PlayerConfig{}, fmt.Errorf("bad data type for config: %T", data["config"])
	}
	cfg := PlayerConfig{}
	if _, ok := m["tileset"]; ok {
		if tileset, ok := m["tileset"].(string); ok {
			cfg.Tileset = tileset
		} else {
			return PlayerConfig{}, fmt.Errorf("bad data type for tileset: %T", m["tileset"])
		}
	}
	return cfg, nil
}

func (s *FireStore) GetPlayer(ctx context.Context, playerID string) (PlayerInfo, error) {
	doc, err := s.players().Doc(playerID).Get(ctx)
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
	if cfg, err := playerConfigFromData(doc.Data()); err == nil {
		pi.Config = cfg
	} else {
		return PlayerInfo{}, fmt.Errorf("could not get config: %v", err)
	}

	return pi, nil
}

func (s *FireStore) GetPlayerByName(ctx context.Context, playerName string) (PlayerInfo, error) {
	iter := s.players().Where("name", "==", playerName).Documents(ctx)
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
	if cfg, err := playerConfigFromData(docs[0].Data()); err == nil {
		pi.Config = cfg
	} else {
		return PlayerInfo{}, fmt.Errorf("could not get config: %v", err)
	}

	return pi, nil
}

func (s *FireStore) RecordPlayerActive(ctx context.Context, code, playerName string, lastActive int64) error {
	_, err := s.games().Doc(code).Collection("active").Doc(playerName).Set(ctx, map[string]any{
		"last_active": lastActive,
	})
	return err
}

func (s *FireStore) PlayerLastActive(ctx context.Context, code, playerName string) (int64, error) {
	doc, err := s.games().Doc(code).Collection("active").Doc(playerName).Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not read: %v", err)
	}
	if lastActive, ok := doc.Data()["last_active"].(int64); ok {
		return lastActive, nil
	}
	return 0, fmt.Errorf("bad data type for last_active: %T", doc.Data()["last_active"])
}

func (s *FireStore) UpdatePlayerConfig(ctx context.Context, playerID string, config PlayerConfig) error {
	clog.Info(ctx, "updating player config", "playerID", playerID)
	_, err := s.players().Doc(playerID).Update(ctx, []firestore.Update{{
		Path:  "config",
		Value: config,
	}})
	return err
}

func (s *FireStore) ReportIssue(ctx context.Context, playerName string, game *Game, summary, whatHappened, whatShouldHappen, errorMessage string) error {
	gameData, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("could not marshal: %v", err)
	}
	_, err = s.issues().NewDoc().Set(ctx, map[string]any{
		"reported_by":      playerName,
		"summary":          summary,
		"whatHappened":     whatHappened,
		"whatShouldHappen": whatShouldHappen,
		"errorMessage":     errorMessage,
		"game_json":        string(gameData),
	})
	return err
}

func (s *FireStore) iterToSummaries(ctx context.Context, iter *firestore.DocumentIterator) ([]GameSummary, error) {
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get all: %v", err)
	}
	summaries := make([]GameSummary, 0, len(docs))

	for _, doc := range docs {
		ctx = clog.WithKeyword(ctx, "code", doc.Ref.ID)
		data := doc.Data()
		scoreboard, ok := data["scoreboard"].(map[string]any)
		if !ok {
			err := fmt.Errorf("bad data type for scoreboard: %T", data["scoreboard"])
			clog.Error(ctx, "bad data type for scoreboard", err)
			continue
		}
		scoreboard64 := make(map[string]int64)
		for name, score := range scoreboard {
			scoreboard64[name] = score.(int64)
		}
		updated, ok := data["updated"].(int64)
		if !ok {
			err := fmt.Errorf("bad data type for updated: %T", data["updated"])
			clog.Error(ctx, "bad data type for updated", err)
			continue
		}
		summary := GameSummary{
			Code:       doc.Ref.ID,
			Scoreboard: scoreboard64,
			Updated:    updated,
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (s *FireStore) queryToSummaries(ctx context.Context, query firestore.Query, updated int64) ([]GameSummary, error) {
	timesThrough := 0
	snaps := query.Snapshots(ctx)
	defer snaps.Stop()
	for {
		snap, err := snaps.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get next snapshot: %w", err)
		}
		summaries, err := s.iterToSummaries(ctx, snap.Documents)
		if err != nil {
			return nil, fmt.Errorf("could not get summaries: %w", err)
		}
		isFresh := updated == 0 // always fresh the first time.
		for _, summary := range summaries {
			if summary.Updated > updated {
				isFresh = true
				break
			}
		}
		if !isFresh && timesThrough == 0 {
			clog.Info(ctx, "not fresh", "timesThrough", timesThrough)
			timesThrough++
			continue
		}
		return summaries, nil
	}
}

func (s *FireStore) ListPickupGames(ctx context.Context, count int, updated int64) ([]GameSummary, error) {
	query := s.scoreboards().
		Where("pickup", "==", true).
		Where("done", "==", false).
		Where("open", "==", true).
		OrderBy("updated", firestore.Desc).
		Limit(count)
	return s.queryToSummaries(ctx, query, updated)
}

func (s *FireStore) ListActiveGames(ctx context.Context, count int, updated int64) ([]GameSummary, error) {
	cutoff := time.Now().Add(-20 * time.Minute).Unix()
	query := s.scoreboards().
		Where("done", "==", false).
		Where("open", "==", false).
		Where("updated", ">", cutoff).
		OrderBy("updated", firestore.Desc).
		Limit(count)
	return s.queryToSummaries(ctx, query, updated)
}

func (s *FireStore) ListRecentGames(ctx context.Context, count int, updated int64) ([]GameSummary, error) {
	query := s.scoreboards().
		Where("done", "==", true).
		Where("has_nonzero_score", "==", true).
		OrderBy("updated", firestore.Desc).
		Limit(count)
	return s.queryToSummaries(ctx, query, updated)
}

func (s *FireStore) waitForUpdatedScoreboards(ctx context.Context, updated int64) (int64, error) {
	query := s.scoreboards().Where("updated", ">", updated).Limit(1)
	snaps := query.Snapshots(ctx)
	defer snaps.Stop()
	for {
		snap, err := snaps.Next()
		if err != nil {
			return 0, fmt.Errorf("could not get next snapshot: %w", err)
		}
		docs, err := snap.Documents.GetAll()
		if err != nil {
			return 0, fmt.Errorf("could not get all: %w", err)
		}
		lastUpdate := updated
		if len(docs) > 0 {
			docUpdated, ok := docs[0].Data()["updated"].(int64)
			if !ok {
				return 0, fmt.Errorf("bad data type for updated: %T", docs[0].Data()["updated"])
			}
			if docUpdated > lastUpdate {
				lastUpdate = docUpdated
			}
			return lastUpdate, nil
		}
	}
}

func (s *FireStore) ListScoreboards(ctx context.Context, count int, updated int64) (ScoreboardSummary, error) {
	lastUpdate, err := s.waitForUpdatedScoreboards(ctx, updated)
	if err != nil {
		return ScoreboardSummary{}, fmt.Errorf("could not wait for updated scoreboards: %w", err)
	}
	active, err := s.ListActiveGames(ctx, count, 0)
	if err != nil {
		return ScoreboardSummary{}, fmt.Errorf("could not get active games: %w", err)
	}
	pickup, err := s.ListPickupGames(ctx, count, 0)
	if err != nil {
		return ScoreboardSummary{}, fmt.Errorf("could not get pickup games: %w", err)
	}
	recent, err := s.ListRecentGames(ctx, count, 0)
	if err != nil {
		return ScoreboardSummary{}, fmt.Errorf("could not get recent games: %w", err)
	}

	lastUpdated := lastUpdate
	for _, summary := range active {
		if summary.Updated > lastUpdated {
			lastUpdated = summary.Updated
		}
	}
	for _, summary := range pickup {
		if summary.Updated > lastUpdated {
			lastUpdated = summary.Updated
		}
	}
	for _, summary := range recent {
		if summary.Updated > lastUpdated {
			lastUpdated = summary.Updated
		}
	}
	// Usually we'll have at least some recent games, but this bootstraps.
	if lastUpdated == 0 {
		lastUpdated = time.Now().Unix()
	}

	return ScoreboardSummary{
		Active:  active,
		Pickup:  pickup,
		Recent:  recent,
		Updated: lastUpdated,
	}, nil
}
