package game

import (
	"errors"

	"math/rand"
)

var (
	ErrGameAlreadyStarted       = errors.New("game already started")
	ErrGameTooManyPlayers       = errors.New("game already has 6 players")
	ErrGameNotEnoughPlayers     = errors.New("not enough players")
	ErrGamePreviousRoundNotDone = errors.New("previous round not done")
	ErrPlayerAlreadyInGame      = errors.New("player already in game")
)

var Colors = []string{"red", "blue", "green"}

type Game struct {
	Version int `json:"version"`

	Code        string    `json:"code"`
	Players     []*Player `json:"players"`
	Rounds      []Round   `json:"rounds"`
	BoardWidth  int       `json:"board_width"`
	BoardHeight int       `json:"board_height"`
}

func NewGame(code string) *Game {
	return &Game{
		Code:        code,
		Version:     0,
		BoardWidth:  10,
		BoardHeight: 11,
	}
}

func (g *Game) AddPlayer(player *Player) error {
	if len(g.Players) >= 6 {
		return ErrGameTooManyPlayers
	}

	if len(g.Rounds) > 0 {
		return ErrGameAlreadyStarted
	}

	for _, p := range g.Players {
		if p.Name == player.Name {
			return ErrPlayerAlreadyInGame
		}
	}

	player.Score = 0
	g.Players = append(g.Players, player)

	return nil
}

func (g *Game) Start() error {
	if len(g.Players) < 2 {
		return ErrGameNotEnoughPlayers
	}

	if len(g.Rounds) > 0 {
		if !g.Rounds[len(g.Rounds)-1].Done() {
			return ErrGamePreviousRoundNotDone
		}
	}

	g.Rounds = append(g.Rounds, Round{
		Turn:      0,
		LaidTiles: []LaidTile{},
	})

	for _, p := range g.Players {
		p.Hand = []Tile{{
			PipsA: rand.Intn(16),
			PipsB: rand.Intn(16),
		}}
	}

	return nil
}

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Hand  []Tile `json:"hand"`
}

type Tile struct {
	PipsA int `json:"pips_a"`
	PipsB int `json:"pips_b"`
}

type Round struct {
	Turn      int        `json:"turn"`
	LaidTiles []LaidTile `json:"laid_tiles"`
}

func (r *Round) Done() bool {
	return false
}

type LaidTile struct {
	Tile        Tile   `json:"tile"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Orientation string `json:"orientation"`
	PlayerName  string `json:"player_name"`
}
