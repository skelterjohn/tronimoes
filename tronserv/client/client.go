package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (c *TronimoesClient) Do(ctx context.Context, method, path string, vin, vout any) error {
	var body io.Reader
	if vin == nil {
		body = http.NoBody
	} else {
		b := &bytes.Buffer{}
		if err := json.NewEncoder(b).Encode(vin); err != nil {
			return fmt.Errorf("could not encode request body: %v", err)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.TronservAddr, path), body)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("could not do request: %s", err)
	}

	if vout != nil {
		if err := json.NewDecoder(resp.Body).Decode(vout); err != nil {
			return fmt.Errorf("could not decode response: %s", err)
		}
	}

	return nil
}

func (c *TronimoesClient) GetPlayer(ctx context.Context, name string) (*game.PlayerInfo, error) {
	var pi game.PlayerInfo
	if err := c.Do(ctx, "GET", "players/jt", nil, &pi); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", pi)

	return nil, nil
}
