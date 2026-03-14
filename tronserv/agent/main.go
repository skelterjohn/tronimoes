package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"

	"github.com/skelterjohn/tronimoes/tronserv/agent/gibbs_planner"
	"github.com/skelterjohn/tronimoes/tronserv/agent/reacts"
	"github.com/skelterjohn/tronimoes/tronserv/agent/types"
	"github.com/skelterjohn/tronimoes/tronserv/client"
	"github.com/skelterjohn/tronimoes/tronserv/game"
)

var (
	tronserv_addr = flag.String("addr", "http://localhost:8080", "host/port for the tronimoes game server")
	name          = flag.String("name", "", "name of the agent")
	gamecode      = flag.String("code", "PICKUP", "code of the game to connect to")
	which         = flag.String("which", "random", "which agent to use: random, gibbs")
	minMoveTime   = flag.Duration("min-move-time", 3*time.Second, "minimum time between moves")
	useGCEToken   = flag.Bool("gce", false, "use the runner's service account to inject access tokens into requests")
	archive       = flag.String("archive", "", "directory to save JSON game state and chosen move per turn; empty = don't save")
	roundOut      = flag.Int("round-out", 0, "targeted player count")
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

func quitFromRoundOut(ctx context.Context, g *game.Game, name string, targetPlayerCount int) bool {
	if targetPlayerCount == 0 {
		return false
	}
	if len(g.Players) <= targetPlayerCount {
		return false
	}
	// a bot needs to quit. am I the one with the highest number?
	highestBotNumber := -1
	for i, p := range g.Players {
		// only bots are allowed multi-word names
		if strings.Contains(p.Name, " ") {
			highestBotNumber = i
		}
	}
	if g.Players[highestBotNumber].Name == name {
		return true
	}
	return false
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

	name := *name

	tc := &client.TronimoesClient{
		TronservAddr: *tronserv_addr,
		Client:       c,
		Name:         name,
	}

	var a types.Agent
	switch *which {
	case "random":
		if name == "" {
			name = CreateName("RC")
		}
		a = RandomChoice{}
	case "gibbs":
		if name == "" {
			name = CreateName("GP")
		}
		gp := &gibbs_planner.GibbsPlanner{
			Name:   name,
			Client: tc,
		}
		gp.SetDefaults()
		a = gp
	default:
		log.Fatalf("Unknown agent: %s", *which)
	}

	tc.Name = name

	log.Printf("Starting %s agent %s, connecting to %s for game %s", *which, name, *tronserv_addr, *gamecode)
	if *archive != "" {
		if err := os.MkdirAll(*archive, 0755); err != nil {
			log.Fatalf("save-dir: %v", err)
		}
	}

	g, err := tc.JoinGame(ctx, *gamecode)
	if err != nil {
		log.Printf("Could not join game: %v", err)
		return
	}

	log.Printf("Joined game %s", g.Code)

	lastUpdateGame := g

	lastMoveTime := time.Now()

	footURL, err := reacts.FindImageURL(ctx, "bot")
	if err != nil {
		log.Printf("Could not get image URL: %v", err)
	} else {
		g, err = tc.ChooseFoot(ctx, footURL)
		if err != nil {
			log.Printf("Could not choose foot: %v", err)
		}
	}

	roundDoneCounter := -1

	for !g.Done {
		if len(g.Rounds) == 0 {
			log.Print("New game beginning")
			if quitFromRoundOut(ctx, g, name, *roundOut) {
				log.Print("Round out reached, quitting to leave room")
				g, err = tc.LeaveOrQuit(ctx)
				if err != nil {
					log.Printf("Could not leave game: %v", err)
				}
				return
			}
		} else if g.Rounds[len(g.Rounds)-1].Done {
			if roundDoneCounter < len(g.Rounds) {
				log.Print("Round done")
				a.CompleteRound(ctx, g)
				roundDoneCounter = len(g.Rounds)
			}
		}

		r := g.CurrentRound(ctx)
		if r == nil || r.Done {
			p := g.GetPlayer(ctx, name)
			if !p.Ready {
				a.Ready(ctx)
				log.Print("Ready to begin a new round.")
				g, err = tc.Start(ctx)
				if err != nil {
					log.Printf("Error starting game: %v", err)
					return
				}
			}
			if g.CurrentRound(ctx) != nil {
				a.Update(ctx, lastUpdateGame, g)
				lastUpdateGame = g
			}
		} else {
			a.Update(ctx, lastUpdateGame, g)
			lastUpdateGame = g
		}
		if len(g.Rounds) > 0 && !g.Rounds[len(g.Rounds)-1].Done {
			if g.Players[g.Turn].Name == name {
				log.Printf("It's my turn")
				log.Printf(" %v", g.Players[g.Turn].Hand)
			} else {
				log.Printf("It's %s's turn", g.Players[g.Turn].Name)
			}
			if g.Players[g.Turn].Name == name {
				p := g.GetPlayer(ctx, name)
				m := a.GetMove(ctx, g, p)
				log.Printf("Move: %+v", m)
				if *archive != "" {
					path := filepath.Join(*archive, fmt.Sprintf("%s_%d.json", g.Code, g.Version))
					blob, err := json.MarshalIndent(struct {
						Game *game.Game `json:"game"`
						Move types.Move `json:"move"`
					}{Game: g, Move: m}, "", "\t")
					if err != nil {
						log.Printf("save marshal: %v", err)
					} else if err := os.WriteFile(path, blob, 0644); err != nil {
						log.Printf("save %s: %v", path, err)
					}
				}
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
				if m.LayTile {
					g, err = tc.LayTile(ctx, &m.LaidTile)
					if err != nil {
						log.Printf("Could not lay tile: %v", err)
						return
					}
					log.Printf("I laid %v", m.LaidTile)
					continue
				}
				if m.PlaceSpacer {
					g, err = tc.LaySpacer(ctx, &m.Spacer)
					if err != nil {
						log.Printf("Could not lay spacer: %v", err)
						return
					}
					log.Printf("I placed a spacer: %v", m.Spacer)
					continue
				}
				log.Println("I did not make a move")
			}
		}

		previousGame := g
		g, err = tc.GetGame(ctx, previousGame.Version)
		for err != nil || g.Version == previousGame.Version {
			if err != nil && err != client.ErrTimeout {
				log.Printf("Game fetch error: %v", err)
				return
			}
			time.Sleep(5 * time.Second)
			g, err = tc.GetGame(ctx, previousGame.Version)
		}

		lastMoveTime = time.Now()

		if len(g.Rounds) == 0 {
			continue
		}

		move, ok := types.InferMove(ctx, previousGame, g)
		if ok {
			pp := previousGame.Players[previousGame.Turn]
			log.Printf("%s played: %s", pp.Name, move)
		}

	}
	a.CompleteGame(ctx, g)
	log.Println("Game over")
}
