package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	"github.com/skelterjohn/tronimoes/tronserv/clog"
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

type GCEMetadataRoundTripper struct {
	Next     http.RoundTripper
	TokenURL string
}

func (a *GCEMetadataRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
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
		if p.Bot {
			highestBotNumber = i
		}
	}
	if g.Players[highestBotNumber].Name == name {
		return true
	}
	return false
}

func areWeAllBots(ctx context.Context, g *game.Game) bool {
	for _, p := range g.Players {
		if !p.Bot {
			return false
		}
	}
	return true
}

func main() {
	ctx := context.Background()
	flag.Parse()

	ctx = clog.WithCloudLoggingOutput(ctx, "tronagent")
	defer clog.CloseCloudLogging(ctx)
	ctx = clog.WithSeverities(ctx, "info", "error")

	c := http.DefaultClient
	if *useGCEToken {
		c = &http.Client{
			Transport: &GCEMetadataRoundTripper{
				Next:     http.DefaultTransport,
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
		clog.Error(ctx, "Unknown agent", nil, "agent", *which)
		os.Exit(1)
	}
	ctx = clog.WithKeyword(ctx, "which", *which)

	tc.Name = name
	ctx = clog.WithKeyword(ctx, "name", name)

	ctx = clog.WithKeyword(ctx, "code", *gamecode)

	clog.Info(ctx, fmt.Sprintf("Starting agent and connecting to game server", "addr", *tronserv_addr))
	if *archive != "" {
		if err := os.MkdirAll(*archive, 0755); err != nil {
			clog.Error(ctx, "Could not create archive directory", err)
			os.Exit(1)
		}
	}

	g, err := tc.JoinGame(ctx, *gamecode)
	if err != nil {
		clog.Error(ctx, "Could not join game", err)
		return
	}

	defer func() {
		if _, err := tc.LeaveOrQuit(ctx); err != nil {
			clog.Error(ctx, "Could not leave game", err)
		}
	}()

	clog.Info(ctx, "Joined game")

	lastUpdateGame := g

	lastMoveTime := time.Now()

	footURL, err := reacts.FindImageURL(ctx, "bot")
	if err != nil {
		clog.Error(ctx, "Could not get image URL", err)
	} else {
		if ng, err := tc.ChooseFoot(ctx, footURL); err != nil {
			clog.Error(ctx, "Could not choose foot", err)
		} else {
			g = ng
		}
	}

	roundDoneCounter := -1

	if len(g.Rounds) == 0 {
		clog.Info(ctx, "New game beginning")
	}

	for !g.Done {
		if len(g.Rounds) == 0 {
			if quitFromRoundOut(ctx, g, name, *roundOut) {
				clog.Info(ctx, "Round out reached, quitting to leave room", "AgentRoundOut", *roundOut)
				return
			}
		} else if g.Rounds[len(g.Rounds)-1].Done {
			if roundDoneCounter < len(g.Rounds) {
				clog.Info(ctx, "Round done")
				a.CompleteRound(ctx, g)
				roundDoneCounter = len(g.Rounds)
			}
		}

		if areWeAllBots(ctx, g) {
			clog.Info(ctx, "All bots, quitting to save $$$$$$$")
			return
		}

		r := g.CurrentRound(ctx)
		if r == nil || r.Done {
			p := g.GetPlayer(ctx, name)
			if p == nil {
				clog.Error(ctx, "Player not found", nil, "player", name)
				return
			}
			if !p.Ready {
				a.Ready(ctx)
				clog.Info(ctx, "Ready to begin a new round.")
				g, err = tc.Start(ctx)
				if err != nil {
					clog.Error(ctx, "Error starting game", err)
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
				clog.Info(ctx, "It's my turn")
				clog.Info(ctx, fmt.Sprintf(" %v", g.Players[g.Turn].Hand))
			} else {
				clog.Info(ctx, fmt.Sprintf("It's %s's turn", g.Players[g.Turn].Name))
			}
			if g.Players[g.Turn].Name == name {
				p := g.GetPlayer(ctx, name)
				m := a.GetMove(ctx, g, p)
				clog.Info(ctx, fmt.Sprintf("Move: %+v", m))
				if *archive != "" {
					path := filepath.Join(*archive, fmt.Sprintf("%s_%d.json", g.Code, g.Version))
					blob, err := json.MarshalIndent(struct {
						Game *game.Game `json:"game"`
						Move types.Move `json:"move"`
					}{Game: g, Move: m}, "", "\t")
					if err != nil {
						clog.Error(ctx, "save marshal", err)
					} else if err := os.WriteFile(path, blob, 0644); err != nil {
						clog.Error(ctx, "Could not write to archive", err, "path", path)
					}
				}
				if time.Since(lastMoveTime) < *minMoveTime {
					// Always wait at least 3 seconds between moves, so
					// as not to confuse the normies.
					time.Sleep(*minMoveTime - time.Since(lastMoveTime))
				}
				lastMoveTime = time.Now()
				if m.Draw {
					if ng, err := tc.Draw(ctx); err != nil {
						clog.Error(ctx, "Could not draw", err)
						continue
					} else {
						g = ng
					}
					clog.Info(ctx, "I just drew")
					continue
				}
				if m.Pass {
					if ng, err := tc.Pass(ctx, m.Selected.X, m.Selected.Y); err != nil {
						clog.Error(ctx, "Could not pass", err)
						continue
					} else {
						g = ng
					}
					clog.Info(ctx, "I passed")
					continue
				}
				if m.LayTile {
					if ng, err := tc.LayTile(ctx, &m.LaidTile); err != nil {
						clog.Error(ctx, "Could not lay tile", err)
						continue
					} else {
						g = ng
					}
					clog.Info(ctx, fmt.Sprintf("I laid %v", m.LaidTile))
					continue
				}
				if m.PlaceSpacer {
					if ng, err := tc.LaySpacer(ctx, &m.Spacer); err != nil {
						clog.Error(ctx, "Could not lay spacer", err)
						continue
					} else {
						g = ng
					}
					clog.Info(ctx, fmt.Sprintf("I placed a spacer: %v", m.Spacer))
					continue
				}
				clog.Info(ctx, "I did not make a move")
			}
		}

		previousGame := g
		g, err = tc.GetGame(ctx, previousGame.Version)
		for err != nil || g.Version == previousGame.Version {
			if err != nil && err != client.ErrTimeout {
				clog.Error(ctx, "Game fetch error", err)
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
			clog.Info(ctx, fmt.Sprintf("%s played: %s", pp.Name, move))
		}

	}
	a.CompleteGame(ctx, g)
	clog.Info(ctx, "Game over")
}
