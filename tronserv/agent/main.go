package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"

	"github.com/skelterjohn/tronimoes/tronserv/client"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

var (
	tronserv_addr = flag.String("addr", "http://localhost:8080", "host/port for the tronimoes game server")
	name          = flag.String("name", "", "name of the agent")
	gamecode      = flag.String("code", "", "code of the game to connect to")
	useGCEToken   = flag.Bool("gce", false, "use the runner's service account to inject access tokens into requests")
)

type AgentRoundTripper struct {
	Next     http.RoundTripper
	TokenURL string
}

func (a *AgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := metadata.GetWithContext(req.Context(), a.TokenURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch id token: %w", err)
	}
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	return a.Next.RoundTrip(req)
}

type Selected struct {
	X int
	Y int
}

type Move struct {
	LaidTile *game.Tile
	Spacer   *game.Spacer
	Draw     bool
	Pass     bool
	Selected Selected
}

type Agent interface {
	Ready()
	Update(g *game.Game)
	GetMove(g *game.Game) Move
}

func main() {
	ctx := context.Background()
	flag.Parse()

	c := http.DefaultClient
	if *useGCEToken {
		c = &http.Client{
			Transport: &AgentRoundTripper{
				Next:     http.DefaultClient.Transport,
				TokenURL: fmt.Sprintf("instance/service-accounts/default/identity?audience=%s", *tronserv_addr),
			},
		}
	}

	tc := client.TronimoesClient{
		TronservAddr: *tronserv_addr,
		Client:       c,
		Name:         *name,
	}

	g, err := tc.JoinGame(ctx, *gamecode)
	if err != nil {
		log.Printf("Could not join game: %v", err)
		return
	}

	a := RandomAgent{}

	a.Ready()

	g, err = tc.Start(ctx)
	if err != nil {
		log.Printf("Error starting game: %v", err)
		return
	}

	for !g.Done {
		if len(g.Rounds) > 0 && !g.Rounds[len(g.Rounds)-1].Done {
			log.Printf("It's %s's turn", g.Players[g.Turn].Name)
			if g.Players[g.Turn].Name == *name {
				p := g.GetPlayer(ctx, *name)
				m := a.GetMove(g, p)
				if m.Draw {
					g, err = tc.Draw(ctx)
					if err != nil {
						log.Printf("Could not draw: %v", err)
						return
					}
					log.Println("drew")
					continue
				}
				if m.Pass {
					g, err = tc.Pass(ctx, m.Selected.X, m.Selected.Y)
					if err != nil {
						log.Printf("Could not pass: %v", err)
						return
					}
					log.Println("passed")
					continue
				}
			} else {
				a.Update(g)
			}
		}

		vers := g.Version
		g, err = tc.GetGame(ctx, vers)
		for g.Version == vers || err == client.ErrTimeout {
			time.Sleep(5 * time.Second)
			g, err = tc.GetGame(ctx, vers)
		}
		if err != nil {
			log.Printf("Could not get game: %v", err)
			return
		}
	}
}
