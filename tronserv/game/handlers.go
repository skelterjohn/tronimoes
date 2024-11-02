package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r chi.Router, s Store) {
	gs := &GameServer{store: s}
	r.Get("/game/{code}", gs.HandleGetGame)
	r.Put("/game/{code}", gs.HandlePutGame)
	r.Post("/game/{code}/start", gs.HandleStartRound)
}

type GameServer struct {
	store Store
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	g, err := s.store.ReadGame(ctx, code)
	if err != nil || g == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// We aleady have something newer.
	if g.Version > version {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if g == nil {
		g = NewGame(code)
	}

	if len(g.Players) >= 4 {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "game already has 4 players")
		return
	}

	player := &Player{}
	if err := json.NewDecoder(r.Body).Decode(player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.AddPlayer(player)
	if err := s.store.WriteGame(ctx, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *GameServer) HandleStartRound(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	g, err := s.store.ReadGame(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := g.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.WriteGame(r.Context(), g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
}
