package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

type TronimoesClient struct {
	TronservAddr string
	Client       *http.Client

	Name string
}

func (c *TronimoesClient) WriteHeaders(req *http.Request) {
	req.Header.Set("X-Player-Name", c.Name)
}

func (c *TronimoesClient) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", c.TronservAddr, path), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	return c.Client.Do(req)
}

func (c *TronimoesClient) GetPlayer(ctx context.Context, name string) (*game.PlayerInfo, error) {
	resp, err := c.Get(ctx, "players/jt")
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err)
	}

	var pi game.PlayerInfo

	if err := json.NewDecoder(resp.Body).Decode(&pi); err != nil {
		return nil, fmt.Errorf("could not decode player: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Printf("Error closing player response: %v", err)
	}

	fmt.Printf("%+v\n", pi)

	return nil, nil
}
