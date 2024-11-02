package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

var (
	addr = flag.String("addr", "0.0.0.0", "address to listen on")
	port = flag.Int("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	game.RegisterHandlers(r, game.NewMemoryStore())

	listenAddr := fmt.Sprintf("%s:%d", *addr, *port)
	log.Printf("Server starting on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatal(err)
	}
}
