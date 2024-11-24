package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
)

var fbApp *firebase.App

func init() {
	ctx := context.Background()
	var err error
	fbApp, err = firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "tronimoes",
	})
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}
}

func writeErr(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func RegisterHandlers(r chi.Router, s Store) {
	gs := &GameServer{store: s}
	r.Get("/game/{code}", gs.HandleGetGame)
	r.Put("/game/{code}", gs.HandlePutGame)
	r.Post("/game/{code}/start", gs.HandleStartRound)
	r.Post("/game/{code}/tile", gs.HandleLayTile)
	r.Post("/game/{code}/spacer", gs.HandleLaySpacer)
	r.Post("/game/{code}/draw", gs.HandleDrawTile)
	r.Post("/game/{code}/pass", gs.HandlePass)
	r.Post("/game/{code}/leave", gs.HandleLeaveOrQuit)
	r.Post("/game/{code}/foot", gs.HandleChickenFoot)
	r.Post("/players", gs.HandleRegisterPlayerName)
	r.Get("/players", gs.HandleGetPlayerName)
}

func RandomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type GameServer struct {
	store Store
}

func (s *GameServer) validateToken(ctx context.Context, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ErrMissingToken
	}
	userID := r.Header.Get("X-Player-Id")
	if userID == "" {
		return ErrMissingUserID
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Verify the Firebase token
	client, err := fbApp.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	decodedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return fmt.Errorf("error verifying ID token: %v", err)
	}

	// Verify that the token's UID matches the X-Player-Id
	if decodedToken.UID != userID {
		return ErrInvalidToken
	}

	return nil
}

func (s *GameServer) getName(ctx context.Context, r *http.Request) (string, error) {
	userID := r.Header.Get("X-Player-Id")
	if userID != "" {
		if err := s.validateToken(ctx, r); err != nil {
			return "", err
		}

		pi, err := s.store.GetPlayer(ctx, userID)
		if err == nil {
			return pi.Name, nil
		}
		return "", err
	}

	tempName := r.Header.Get("X-Player-Name")
	_, err := s.store.GetPlayerByName(ctx, tempName)
	if err == ErrNoRegisteredPlayer {
		// anonymous play is ok with unregistered names.
		return tempName, nil
	}
	if err == nil {
		return "", ErrNotYourPlayer
	}
	return "", err
}

func (s *GameServer) encodeFilteredGame(ctx context.Context, w http.ResponseWriter, name string, g *Game) {
	for _, p := range g.Players {
		if p.Name == name {
			continue
		}
		// Hide the hands of other players, though we still send the tile counts.
		for _, t := range p.Hand {
			t.PipsA = 0
			t.PipsB = 0
		}
	}
	// Hide the bag from everyone.
	for _, t := range g.Bag {
		t.PipsA = 0
		t.PipsB = 0
	}

	// Add legal moves for this player to see.
	r := g.CurrentRound(ctx)
	if len(g.Players) > g.Turn {
		p := g.Players[g.Turn]
		if !g.Done && r != nil && p.Name == name {
			r.FindHints(ctx, g, name, p)
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		log.Printf("Error encoding game %q: %v", g.Code, err)
		writeErr(w, err, http.StatusInternalServerError)
	}
}

func (s *GameServer) HandleLeaveOrQuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if !g.LeaveOrQuit(ctx, name) {
		log.Printf("Player %q cannot leave game %q", name, code)
		writeErr(w, ErrNotYourGame, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleDrawTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	g.CheckForDupes(ctx, "before draw")

	player := g.Players[g.Turn]
	if player.Name != name {
		log.Printf("Player %q is not in turn for game %q", name, code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		log.Printf("Player %q tried to play game %q but it isn't started", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		log.Printf("Player %q tried to play game %q but the round is done", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	if !g.DrawTile(ctx, name) {
		log.Printf("Player %q tried to play game %q but it isn't started", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	g.CheckForDupes(ctx, "after draw")
	s.encodeFilteredGame(ctx, w, name, g)
}

type ChickenFoot struct {
	SelectedX int `json:"selected_x"`
	SelectedY int `json:"selected_y"`
}

func (s *GameServer) HandlePass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	chickenFoot := &ChickenFoot{}
	if err := json.NewDecoder(r.Body).Decode(chickenFoot); err != nil {
		log.Printf("Error decoding chicken-foot pacement %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "before pass")

	player := g.Players[g.Turn]
	if player.Name != name {
		log.Printf("Player %q is not in turn for game %q", name, code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		log.Printf("Player %q tried to play game %q but it isn't started", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		log.Printf("Player %q tried to play game %q but the round is done", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	if err := g.Pass(ctx, name, chickenFoot.SelectedX, chickenFoot.SelectedY); err != nil {
		log.Printf("Could not pass: %v", err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	g.CheckForDupes(ctx, "after pass")

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleLayTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "before lay")

	player := g.Players[g.Turn]
	if player.Name != name {
		log.Printf("Player %q is not in turn for game %q", name, code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		log.Printf("Player %q tried to play game %q but it isn't started", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		log.Printf("Player %q tried to play game %q but the round is done", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	lt := &LaidTile{}
	if err := json.NewDecoder(r.Body).Decode(lt); err != nil {
		log.Printf("Error decoding tile for %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	if lt.Tile == nil {
		log.Printf("No tile provided for %q / %q", name, code)
		writeErr(w, ErrNoTile, http.StatusBadRequest)
		return
	}

	lt.PlayerName = player.Name

	if err := g.LayTile(ctx, name, lt); err != nil {
		log.Printf("Error laying tile for %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "after lay")

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleLaySpacer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "spacer-read")

	player := g.Players[g.Turn]
	if player.Name != name {
		log.Printf("Player %q is not in turn for game %q", name, code)
		writeErr(w, ErrNotYourTurn, http.StatusBadRequest)
		return
	}

	if len(g.Rounds) == 0 {
		log.Printf("Player %q tried to play game %q but it isn't started", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		log.Printf("Player %q tried to play game %q but the round is done", name, code)
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	sp := &Spacer{}
	if err := json.NewDecoder(r.Body).Decode(sp); err != nil {
		log.Printf("Error decoding spacerfor %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := g.LaySpacer(ctx, name, sp); err != nil {
		log.Printf("Error laying spacer for %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "spacer-write")
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Four minutes before timing out. Client must re-request within one minute
	// or risk being booted.
	ctx, cancel := context.WithTimeout(ctx, 200*time.Second)
	defer cancel()

	code := chi.URLParam(r, "code")
	versionStr := r.URL.Query().Get("version")
	var version int64
	if versionStr != "" {
		var err error
		version, err = strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			log.Printf("Error parsing version %q: %v", versionStr, err)
			writeErr(w, err, http.StatusBadRequest)
			return
		}
	}
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	if err := s.store.RecordPlayerActive(ctx, code, name, time.Now().Unix()); err != nil {
		log.Printf("Error setting player active for %q / %q: %v", name, code, err)
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	// if the game hasn't begun, cull inactive players
	if len(g.Rounds) > 0 {
		now := time.Now().Unix()
		anyBooted := false
		for _, p := range g.Players {
			if p.Name == name {
				continue
			}
			lastActive, err := s.store.PlayerLastActive(ctx, code, p.Name)
			if err != nil {
				log.Printf("Error getting last active for %q / %q: %v", p.Name, code, err)
				continue
			}
			idleSeconds := now - lastActive
			if idleSeconds > 300 {
				log.Printf("last active for %q / %q: %d (%d seconds ago)", p.Name, code, lastActive, idleSeconds)
				if !g.LeaveOrQuit(ctx, p.Name) {
					log.Printf("Could not boot inactive player %s from %q", p.Name, code)
				} else {
					anyBooted = true
				}
			}
		}
		if anyBooted {
			log.Printf("Booted players from %q %v", code, g.Players)
			if err := s.store.WriteGame(ctx, g); err != nil {
				log.Printf("Could not store game after booting players: %v", err)
			}
		}
	}

	// We aleady have something newer.
	if g.Version > version {
		g.CheckForDupes(ctx, "get")
		s.encodeFilteredGame(ctx, w, name, g)
		return
	}

	// Otherwise, wait for am update.
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == context.DeadlineExceeded {
			writeErr(w, err, http.StatusRequestTimeout)
			return
		}
		log.Printf("%s broke connection for %q: %v", name, code, err)
		if err := s.store.RecordPlayerActive(ctx, code, name, 0); err != nil {
			log.Printf("Error setting player inactive for %q / %q: %v", name, code, err)
		}
		if !g.LeaveOrQuit(ctx, name) {
			log.Printf("Could not boot %s from %q when the connection broke", name, code)
		}
		return
	case g := <-s.store.WatchGame(ctx, code, version):
		g.CheckForDupes(ctx, "watch")
		s.encodeFilteredGame(ctx, w, name, g)
	}
}

func (s *GameServer) HandlePutGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	if code == "<>" {
		code = "PICKUP"
	}

	g, err := s.store.FindGameAlreadyPlaying(ctx, code, name)
	if err != nil && err != ErrNoSuchGame {
		log.Printf("Error reading game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if g == nil {
		g, err = s.store.FindOpenGame(ctx, code)
		if err != nil && err != ErrNoSuchGame {
			log.Printf("Error reading game %q: %v", code, err)
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
	}

	if g == nil {
		code = fmt.Sprintf("%s-%s", code, RandomString(6))
		g = NewGame(ctx, code)
	}

	inGame := false
	for _, p := range g.Players {
		if p.Name == name {
			inGame = true
			log.Printf("Player %q already in game %q", name, code)
		}
	}

	player := &Player{}
	if err := json.NewDecoder(r.Body).Decode(player); err != nil {
		log.Printf("Error decoding player for %q / %q", name, code)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if player.Name != name {
		log.Printf("Header name %q does not match payload name %q", name, player.Name)
		writeErr(w, ErrNotYou, http.StatusForbidden)
		return
	}

	if !inGame {
		if err := g.AddPlayer(ctx, player); err != nil {
			log.Printf("Error adding player %q to game %q: %v", name, code, err)
			if err == ErrGameTooManyPlayers {
				writeErr(w, err, http.StatusUnprocessableEntity)
				return
			}
			if err == ErrGameAlreadyStarted {
				writeErr(w, err, http.StatusUnprocessableEntity)
				return
			}
			if err == ErrPlayerAlreadyInGame {
				writeErr(w, err, http.StatusConflict)
				return
			}
			writeErr(w, err, http.StatusInternalServerError)
			return
		}

		if err := s.store.WriteGame(ctx, g); err != nil {
			log.Printf("Error writing game %q: %v", code, err)
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
	}

	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleStartRound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "start-read")
	// Only the first player can start the round.
	if g.Players[0].Name != name {
		log.Printf("In game %q, player %q tried to start game for %q", code, name, g.Players[0].Name)
		writeErr(w, ErrNotYourGame, http.StatusForbidden)
		return
	}

	if err := g.Start(ctx); err != nil {
		log.Printf("Error starting round for game %q: %v", code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "start-write")
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleChickenFoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(ctx, r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		log.Printf("Error reading game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "foot-read")

	reqBody := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Error decoding chickenfoot for %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	url, ok := reqBody["url"]
	if !ok {
		log.Printf("No url provided for %q / %q", name, code)
		writeErr(w, ErrNoURL, http.StatusBadRequest)
		return
	}

	player := g.GetPlayer(ctx, name)
	if player == nil {
		log.Printf("Player %q not found in game %q", name, code)
		writeErr(w, ErrPlayerNotFound, http.StatusNotFound)
		return
	}

	player.ChickenFootURL = url

	if err := s.store.WriteGame(ctx, g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	g.CheckForDupes(ctx, "foot-write")
	s.encodeFilteredGame(ctx, w, name, g)
}

func (s *GameServer) HandleRegisterPlayerName(w http.ResponseWriter, r *http.Request) {
	playerID := r.Header.Get("X-Player-ID")

	pi := &PlayerInfo{}
	if err := json.NewDecoder(r.Body).Decode(pi); err != nil {
		log.Printf("Error decoding player info: %v", err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	pi.Id = playerID

	if rpi, err := s.store.GetPlayerByName(r.Context(), pi.Name); err == nil {
		if rpi.Id != playerID {
			log.Printf("Player %q already registered to %q", pi.Name, rpi.Id)
			writeErr(w, ErrPlayerAlreadyRegistered, http.StatusConflict)
			return
		}
	}

	if playerID != "" {
		if err := s.store.RegisterPlayerName(r.Context(), playerID, pi.Name); err != nil {
			log.Printf("Error registering player %q: %v", pi.Name, err)
			writeErr(w, err, http.StatusBadRequest)
			return
		}
		log.Printf("Registered player %q for %s", pi, playerID)
	} else {
		log.Printf("Anonymous player %q", pi.Name)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pi)
}

func (s *GameServer) HandleGetPlayerName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pi, err := s.store.GetPlayer(ctx, r.Header.Get("X-Player-ID"))
	if err != nil {
		if err == ErrNoRegisteredPlayer {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		log.Printf("Error getting player name for %q: %v", r.Header.Get("X-Player-ID"), err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pi)
}
