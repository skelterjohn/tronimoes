package game

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var (
	ErrRoundNotStarted = errors.New("round not started")
)

func RegisterHandlers(r chi.Router, s Store) {
	gs := &GameServer{store: s}
	r.Get("/game/{code}", gs.HandleGetGame)
	r.Put("/game/{code}", gs.HandlePutGame)
	r.Post("/game/{code}/start", gs.HandleStartRound)
	r.Post("/game/{code}/tile", gs.HandleLayTile)
}

type GameServer struct {
	store Store
}

func writeErr(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func (s *GameServer) HandleLayTile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
		if err == ErrNoSuchGame {
			writeErr(w, err, http.StatusNotFound)
			return
		}
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if len(g.Rounds) == 0 {
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}
	round := g.Rounds[len(g.Rounds)-1]
	if round.Done {
		writeErr(w, ErrRoundNotStarted, http.StatusBadRequest)
		return
	}

	lt := &LaidTile{}
	if err := json.NewDecoder(r.Body).Decode(lt); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := g.LayTile(*lt); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
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
			writeErr(w, err, http.StatusBadRequest)
			return
		}
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil {
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
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
		return
	}

	// Otherwise, wait for am update.
	select {
	case <-ctx.Done():
		return
	case game := <-s.store.WatchGame(ctx, code):
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(game); err != nil {
			log.Printf("Error encoding game: %v", err)
			return
		}
	}
}

func (s *GameServer) HandlePutGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")

	g, err := s.store.ReadGame(ctx, code)
	if err != nil && err != ErrNoSuchGame {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if g == nil {
		g = NewGame(code)
	}

	player := &Player{}
	if err := json.NewDecoder(r.Body).Decode(player); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := g.AddPlayer(player); err != nil {
		if err == ErrGameTooManyPlayers {
			writeErr(w, err, http.StatusConflict)
			return
		}
		if err == ErrGameAlreadyStarted {
			writeErr(w, err, http.StatusConflict)
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
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *GameServer) HandleStartRound(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	g, err := s.store.ReadGame(r.Context(), code)
	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	if err := g.Start(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
}
