package queue

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r chi.Router) {
	r.Post("/join/{code}", HandleJoin)
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	fmt.Fprintf(w, "Hello, world! #%s", code)
}
