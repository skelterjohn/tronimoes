package game

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r chi.Router) {
	r.Get("/game/{code}", HandleGetGame)
}

func HandleGetGame(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	g := Game{Code: code}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
