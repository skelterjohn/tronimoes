package game

import (
	"fmt"
	"log"

	"math/rand"
)

var Colors = []string{"red", "blue", "green"}

type Game struct {
	Version int `json:"version"`

	Code        string    `json:"code"`
	Players     []*Player `json:"players"`
	Turn        int       `json:"turn"`
	Rounds      []*Round  `json:"rounds"`
	Bag         []*Tile   `json:"bag"`
	BoardWidth  int       `json:"board_width"`
	BoardHeight int       `json:"board_height"`
	MaxPips     int       `json:"max_pips"`
}

func NewGame(code string) *Game {
	return &Game{
		Code:        code,
		Version:     0,
		BoardWidth:  10,
		BoardHeight: 11,
		MaxPips:     16,
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
		if !g.Rounds[len(g.Rounds)-1].Done {
			return ErrGamePreviousRoundNotDone
		}
	}

	g.Rounds = append(g.Rounds, &Round{
		Turn:      0,
		LaidTiles: []*LaidTile{},
	})

	// Fill the bag with tiles.
	for a := 0; a < g.MaxPips; a++ {
		for b := a; b < g.MaxPips; b++ {
			g.Bag = append(g.Bag, &Tile{
				PipsA: a,
				PipsB: b,
			})
		}
	}
	// And shuffle it.
	rand.Shuffle(len(g.Bag), func(i, j int) {
		g.Bag[i], g.Bag[j] = g.Bag[j], g.Bag[i]
	})

	for _, p := range g.Players {
		for i := 0; i < 7; i++ {
			g.DrawTile(p.Name)
		}
	}

	return nil
}

func (g *Game) DrawTile(name string) bool {
	if len(g.Bag) == 0 {
		return false
	}

	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	if player == nil {
		return false
	}

	player.Hand = append(player.Hand, g.Bag[0])
	log.Printf("%s drew %v", name, g.Bag[0])
	g.Bag = g.Bag[1:]
	return true
}

func (g *Game) GetPlayer(name string) *Player {
	for _, p := range g.Players {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (g *Game) LayTile(tile *LaidTile) error {
	player := g.GetPlayer(tile.PlayerName)
	if player == nil {
		return ErrPlayerNotFound
	}

	newHand := []*Tile{}
	foundTile := false
	for _, t := range player.Hand {
		if t.PipsA != tile.Tile.PipsA || t.PipsB != tile.Tile.PipsB {
			newHand = append(newHand, t)
			continue
		}
		foundTile = true
	}
	if !foundTile {
		return ErrTileNotFound
	}
	player.Hand = newHand

	round := g.Rounds[len(g.Rounds)-1]
	if err := round.LayTile(tile); err != nil {
		return fmt.Errorf("laying tile: %w", err)
	}
	g.Turn = (g.Turn + 1) % len(g.Players)
	return nil
}

type Player struct {
	Name  string  `json:"name"`
	Score int     `json:"score"`
	Hand  []*Tile `json:"hand"`
}

type Tile struct {
	PipsA int `json:"pips_a"`
	PipsB int `json:"pips_b"`
}

type LaidTile struct {
	Tile        *Tile  `json:"tile"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Orientation string `json:"orientation"`
	PlayerName  string `json:"player_name"`
}

type Round struct {
	Turn      int         `json:"turn"`
	LaidTiles []*LaidTile `json:"laid_tiles"`
	Done      bool        `json:"done"`
}

func (r *Round) LayTile(tile *LaidTile) error {
	if r.Done {
		return ErrRoundAlreadyDone
	}
	r.LaidTiles = append(r.LaidTiles, tile)
	return nil
}
