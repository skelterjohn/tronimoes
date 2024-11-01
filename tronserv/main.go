package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	game.RegisterHandlers(r)

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOSTNAME"), os.Getenv("PORT"))
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
