package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/skelterjohn/tronimoes/tronserv/clog"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

// defaultServiceAccountEmail returns the email of the default service account
// for the current GCE/Cloud Run instance via the metadata server. Returns empty
// string and nil when not running on GCE or when the metadata is unavailable.
func defaultServiceAccountEmail(ctx context.Context) string {
	email, err := metadata.GetWithContext(ctx, "instance/service-accounts/default/email")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(email)
}

var (
	addr           = flag.String("addr", "0.0.0.0", "address to listen on")
	port           = flag.Int("port", 8080, "port to listen on")
	env            = flag.String("env", "", "firestore env (unset to use MemoryStore)")
	noCors         = flag.Bool("no-cors", false, "disable cors")
	agentSpawner   = flag.String("agent-spawner", "local", "agent spawner to use: local, gcr, gcr-dev")
	checkBotTokens = flag.Bool("check-bot-tokens", false, "check bot tokens for reserved names")
	dev            = flag.Bool("dev", false, "run in development mode")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	r := chi.NewRouter()

	ctx = clog.WithSeverities(ctx, clog.INFO, clog.ERROR)
	if *dev {
		ctx = clog.WithSeverities(ctx, clog.DEBUG)
		ctx = clog.WithTextOutput(ctx, os.Stdout)
		r.Use(clog.ChiLoggerDev)
		clog.Info(ctx, "Running in development mode")
	} else {
		ctx = clog.WithStructuredOutput(ctx, os.Stdout)
		r.Use(clog.ChiLogger)
		clog.Info(ctx, "Running in production mode")
	}

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
			clog.Info(ctx, "CORS rejected origin", "origin", origin)
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
			clog.Fatal(ctx, "Could not connect to firestore", err)
		}
	}

	var spawner game.AgentSpawner
	switch *agentSpawner {
	case "":
		spawner = nil
	case "local":
		spawner = game.LocalAgentSpawner{}
	case "gcr-dev":
		gcr := &game.GCRAgentSpawner{
			ProjectID: "tronimoes",
			Region:    "us-east4",
			Code:      "BDIKED-BNHJXU",
		}
		if err := gcr.Initialize(ctx); err != nil {
			clog.Error(ctx, "Could not infer GCR agent spawner config", err)
			spawner = nil
		} else {
			spawner = gcr
		}
	case "gcr":
		gcr := &game.GCRAgentSpawner{}
		if err := gcr.Initialize(ctx); err != nil {
			clog.Error(ctx, "Could not infer GCR agent spawner config", err)
			spawner = nil
		} else {
			spawner = gcr
		}
	default:
		clog.Fatal(ctx, "Unknown agent spawner", nil, "agentSpawner", *agentSpawner)
	}

	allowedBotSAs := []string{}
	if *checkBotTokens {
		if sa := defaultServiceAccountEmail(ctx); sa != "" {
			allowedBotSAs = []string{sa}
			clog.Info(ctx, "Bot token check: allowing service account", "serviceAccount", sa)
		}
	}

	gs := &game.GameServer{
		Store:                     store,
		AgentSpawner:              spawner,
		CheckBotTokens:            *checkBotTokens,
		AllowedBotServiceAccounts: allowedBotSAs,
	}
	game.RegisterHandlers(r, gs)

	listenAddr := fmt.Sprintf("%s:%d", *addr, *port)
	clog.Info(ctx, "Server starting", "listenAddr", listenAddr)
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		clog.Fatal(ctx, "Server failed", err)
	}
}
