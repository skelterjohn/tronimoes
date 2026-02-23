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
	which         = flag.String("which", "random", "which agent to use: random, gibbs")
	minMoveTime   = flag.Duration("min-move-time", 3*time.Second, "minimum time between moves")
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

type Move struct {
	LaidTile *game.LaidTile
	Spacer   *game.Spacer
	Draw     bool
	Pass     bool
	Selected game.Coord
}

type Agent interface {
	Ready(ctx context.Context)
	Update(ctx context.Context, g *game.Game)
	GetMove(ctx context.Context, g *game.Game, p *game.Player) Move
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

	log.Printf("Starting %s agent %s, connecting to %s for game %s", *which, *name, *tronserv_addr, *gamecode)

	g, err := tc.JoinGame(ctx, *gamecode)
	if err != nil {
		log.Printf("Could not join game: %v", err)
		return
	}

	lastMoveTime := time.Now()

	var a Agent
	switch *which {
	case "random":
		a = RandomChoice{}
	case "gibbs":
	default:
		log.Fatalf("Unknown agent: %s", *which)
	}

	for !g.Done {
		if len(g.Rounds) == 0 {
			log.Print("New game beginning")
		} else if g.Rounds[len(g.Rounds)-1].Done {
			log.Print("Round done")
		}

		r := g.CurrentRound(ctx)
		if r == nil || r.Done {
			a.Ready(ctx)

			log.Print("Ready to begin a new round.")

			g, err = tc.Start(ctx)
			if err != nil {
				log.Printf("Error starting game: %v", err)
				return
			}
		}
		if len(g.Rounds) > 0 && !g.Rounds[len(g.Rounds)-1].Done {
			if g.Players[g.Turn].Name == *name {
				log.Printf("It's my turn")
			} else {
				log.Printf("It's %s's turn", g.Players[g.Turn].Name)
			}
			if g.Players[g.Turn].Name == *name {
				p := g.GetPlayer(ctx, *name)
				m := a.GetMove(ctx, g, p)
				if time.Since(lastMoveTime) < *minMoveTime {
					// Always wait at least 3 seconds between moves, so
					// as not to confuse the normies.
					time.Sleep(*minMoveTime - time.Since(lastMoveTime))
				}
				lastMoveTime = time.Now()
				if m.Draw {
					g, err = tc.Draw(ctx)
					if err != nil {
						log.Printf("Could not draw: %v", err)
						return
					}
					log.Println("I just drew")
					continue
				}
				if m.Pass {
					g, err = tc.Pass(ctx, m.Selected.X, m.Selected.Y)
					if err != nil {
						log.Printf("Could not pass: %v", err)
						return
					}
					log.Println("I passed")
					continue
				}
				if m.LaidTile != nil {
					g, err = tc.LayTile(ctx, m.LaidTile)
					if err != nil {
						log.Printf("Could not lay tile: %v", err)
						return
					}
					log.Printf("I laid %v", m.LaidTile)
					continue
				}
				if m.Spacer != nil {
					g, err = tc.LaySpacer(ctx, m.Spacer)
					if err != nil {
						log.Printf("Could not lay spacer: %v", err)
						return
					}
					log.Printf("I placed a spacer: %v", m.Spacer)
					continue
				}
				log.Println("I did not make a move")
			} else {
				a.Update(ctx, g)
			}
		}

		previousGame := g
		g, err = tc.GetGame(ctx, previousGame.Version)
		for g.Version == previousGame.Version || err == client.ErrTimeout {
			time.Sleep(5 * time.Second)
			g, err = tc.GetGame(ctx, previousGame.Version)
		}
		if err != nil {
			log.Printf("Could not get game: %v", err)
			return
		}

		lastMoveTime = time.Now()

		currentRound := g.Rounds[len(g.Rounds)-1]
		previousCurrentRound := previousGame.CurrentRound(ctx)
		if currentRound == nil || previousCurrentRound == nil {
			continue
		}
		lastPlayer := g.Players[previousGame.Turn]
		knownPlay := false
		if len(currentRound.LaidTiles) > len(previousCurrentRound.LaidTiles) {
			lastTile := currentRound.LaidTiles[len(currentRound.LaidTiles)-1]
			log.Printf("%s laid %s", lastPlayer.Name, lastTile)
			knownPlay = true
		}
		if currentRound.Spacer != nil {
			log.Printf("%s placed a spacer: %s", lastPlayer.Name, currentRound.Spacer)
			knownPlay = true
		}
		for _, p := range g.Players {
			if p.JustDrew {
				log.Printf("%s just drew", p.Name)
				knownPlay = true
			}
		}
		if !knownPlay && previousGame.Turn != g.Turn {
			log.Printf("%s passed", lastPlayer.Name)
		}
	}
	log.Println("Game over")
}
