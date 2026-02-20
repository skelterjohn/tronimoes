package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

var ErrTimeout = errors.New("request timed out")

type TronimoesClient struct {
	TronservAddr string
	Client       *http.Client

	Name string
	Code string
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
		body = b
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.TronservAddr, path), body)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}

	c.WriteHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("could not do request: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		d, _ := io.ReadAll(resp.Body)
		log.Printf("Got an error with the request: %s", d)
		if resp.StatusCode == http.StatusNotFound {
			return game.ErrNoSuchGame
		}
		if resp.StatusCode == http.StatusRequestTimeout {
			return ErrTimeout
		}
		return fmt.Errorf("request had status code %d", resp.StatusCode)
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

func (c *TronimoesClient) JoinGame(ctx context.Context, code string) (*game.Game, error) {
	p := game.Player{
		Name: c.Name,
	}
	var g game.Game

	if err := c.Do(ctx, "PUT", fmt.Sprintf("game/%s", code), p, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)

	c.Code = g.Code

	return &g, nil
}

func (c *TronimoesClient) GetGame(ctx context.Context, version int64) (*game.Game, error) {
	var g game.Game

	if err := c.Do(ctx, "GET", fmt.Sprintf("game/%s", c.Code), nil, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)

	return &g, nil
}

func (c *TronimoesClient) Start(ctx context.Context) (*game.Game, error) {
	var g game.Game

	if err := c.Do(ctx, "POST", fmt.Sprintf("game/%s/start", c.Code), nil, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)

	return &g, nil
}

func (c *TronimoesClient) Draw(ctx context.Context) (*game.Game, error) {
	var g game.Game

	if err := c.Do(ctx, "POST", fmt.Sprintf("game/%s/draw", c.Code), nil, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)

	return &g, nil
}

func (c *TronimoesClient) Pass(ctx context.Context, x, y int) (*game.Game, error) {
	var g game.Game

	sel := map[string]int{
		"selected_x": x,
		"selected_y": y,
	}

	if err := c.Do(ctx, "POST", fmt.Sprintf("game/%s/pass", c.Code), sel, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)

	return &g, nil
}

func (c *TronimoesClient) LayTile(ctx context.Context, lt *game.LaidTile) (*game.Game, error) {
	var g game.Game
	if err := c.Do(ctx, "POST", fmt.Sprintf("game/%s/tile", c.Code), lt, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)

	return &g, nil
}

func (c *TronimoesClient) LaySpacer(ctx context.Context, sp *game.Spacer) (*game.Game, error) {
	var g game.Game
	if err := c.Do(ctx, "POST", fmt.Sprintf("game/%s/spacer", c.Code), sp, &g); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", g)
	return &g, nil
}
