package game

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
	r.Post("/game/{code}/draw", gs.HandleDrawTile)
	r.Post("/game/{code}/pass", gs.HandlePass)
	r.Post("/game/{code}/leave", gs.HandleLeaveOrQuit)
}

type GameServer struct {
	store Store
}

func (s *GameServer) getName(r *http.Request) (string, error) {
	// TODO: validate bearer token
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", ErrMissingToken
	}
	return r.Header.Get("X-Player-Name"), nil
}

func (s *GameServer) encodeFilteredGame(w http.ResponseWriter, name string, g *Game) {
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		log.Printf("Error encoding game %q: %v", g.Code, err)
		writeErr(w, err, http.StatusInternalServerError)
	}
}

func (s *GameServer) HandleLeaveOrQuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(r)
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

	if !g.LeaveOrQuit(name) {
		log.Printf("Player %q cannot leave game %q", name, code)
		writeErr(w, ErrNotYourGame, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(w, name, g)
}

func (s *GameServer) HandleDrawTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(r)
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

	if !g.DrawTile(player.Name) {
		log.Print("Could not draw a tile")
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(w, name, g)
}

func (s *GameServer) HandlePass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(r)
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

	if !g.Pass(player.Name) {
		log.Print("Could not pass")
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(w, name, g)
}

func (s *GameServer) HandleLayTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(r)
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

	if err := g.LayTile(lt); err != nil {
		log.Printf("Error laying tile for %q / %q: %v", name, code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(w, name, g)
}

func (s *GameServer) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	versionStr := r.URL.Query().Get("version")
	version := 0
	if versionStr != "" {
		var err error
		version, err = strconv.Atoi(versionStr)
		if err != nil {
			log.Printf("Error parsing version %q: %v", versionStr, err)
			writeErr(w, err, http.StatusBadRequest)
			return
		}
	}
	name, err := s.getName(r)
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

	// We aleady have something newer.
	if g.Version > version {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(g); err != nil {
			log.Printf("Error encoding game %q: %v", code, err)
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
		return
	}

	// Otherwise, wait for am update.
	select {
	case <-ctx.Done():
		log.Printf("%s broke connection for %q", name, code)
		return
	case g := <-s.store.WatchGame(ctx, code):
		s.encodeFilteredGame(w, name, g)
	}
}

func (s *GameServer) HandlePutGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(r)
	if err != nil {
		log.Printf("Error getting name: %v", err)
		writeErr(w, err, http.StatusForbidden)
		return
	}

	if code == "<>" {
		code, err = GetPickupCode(name)
		if err != nil {
			log.Printf("Error getting pickup code for %q: %v", name, err)
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil && err != ErrNoSuchGame {
		log.Printf("Error reading game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if g == nil {
		g = NewGame(code)
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
		if err := g.AddPlayer(player); err != nil {
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

	s.encodeFilteredGame(w, name, g)
}

func (s *GameServer) HandleStartRound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	name, err := s.getName(r)
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

	// Only the first player can start the round.
	if g.Players[0].Name != name {
		log.Printf("In game %q, player %q tried to start game for %q", code, name, g.Players[0].Name)
		writeErr(w, ErrNotYourGame, http.StatusForbidden)
		return
	}

	if err := g.Start(); err != nil {
		log.Printf("Error starting round for game %q: %v", code, err)
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		log.Printf("Error writing game %q: %v", code, err)
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	s.encodeFilteredGame(w, name, g)
}
