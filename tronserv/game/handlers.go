package game

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r chi.Router, s Store) {
	gs := &GameServer{store: s}
	r.Get("/game/{code}", gs.HandleGetGame)
	r.Put("/game/{code}", gs.HandlePutGame)
}

type GameServer struct {
	store Store
}

func (s *GameServer) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	g, err := s.store.ReadGame(r.Context(), code)
	if err != nil || g == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		g = &Game{Code: code}
	}

	if len(g.Players) >= 4 {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "game already has 4 players")
		return
	}

	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.Players = append(g.Players, player)
	if err := s.store.WriteGame(ctx, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
