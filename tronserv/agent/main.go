package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"log"

	"cloud.google.com/go/compute/metadata"
)

var (
	tronserv_addr = flag.String("tronserv", "http://localhost:8080", "host/port for the tronimoes game server")
	name          = flag.String("name", "", "name of the agent")
	gamecode      = flag.String("gamecode", "", "code of the game to connect to")
	useGCEToken   = flag.Bool("use_gce_token", false, "use the runner's service account to inject access tokens into requests")
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

func main() {
	ctx := context.Background()
	flag.Parse()

	client := http.DefaultClient
	if *useGCEToken {
		client = &http.Client{
			Transport: &AgentRoundTripper{
				Next:     http.DefaultClient.Transport,
				TokenURL: fmt.Sprintf("instance/service-accounts/default/identity?audience=%s", *tronserv_addr),
			},
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/players/jt", *tronserv_addr), http.NoBody)
	if err != nil {
		log.Printf("Could not create request: %s", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Could not create request: %s", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("%s\n", data)
}
