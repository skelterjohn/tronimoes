package main

import (
	"context"
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
	addr         = flag.String("addr", "0.0.0.0", "address to listen on")
	port         = flag.Int("port", 8080, "port to listen on")
	env          = flag.String("env", "", "firestore env (unset to use MemoryStore)")
	noCors       = flag.Bool("no-cors", false, "disable cors")
	agentSpawner = flag.String("agent-spawner", "local", "agent spawner to use: local, gce")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	allowedOriginsList := []string{"http://localhost:3000", "https://tronapp-1010961884428.us-east4.run.app", "https://tronimoes.com"}
	allowedOrigins := make(map[string]bool)
	for _, o := range allowedOriginsList {
		allowedOrigins[o] = true
	}
	allowAnyOrigin := *noCors
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			if allowAnyOrigin {
				return true
			}
			if allowedOrigins[origin] {
				return true
			}
			log.Printf("CORS rejected origin: %q (allowed: %v)", origin, allowedOriginsList)
			return false
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
			"X-Player-Name", "Authorization",
			"X-Player-ID",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	var store game.Store
	if *env == "" {
		store = game.NewMemoryStore()
	} else {
		var err error
		store, err = game.NewFirestore(ctx, "tronimoes", *env)
		if err != nil {
			log.Fatalf("Could not connect to firestore: %v", err)
		}
	}

	var spawner game.AgentSpawner
	switch *agentSpawner {
	case "":
		spawner = nil
	case "local":
		spawner = game.LocalAgentSpawner{}
	default:
		log.Fatalf("Unknown agent spawner: %s", *agentSpawner)
	}

	gs := &game.GameServer{
		Store:        store,
		AgentSpawner: spawner,
	}
	game.RegisterHandlers(r, gs)

	listenAddr := fmt.Sprintf("%s:%d", *addr, *port)
	log.Printf("Server starting on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatal(err)
	}
}
