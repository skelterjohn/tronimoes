package game

import (
	"fmt"
	"log"

	"math/rand"
)

var Colors = []string{"red", "blue", "green"}

type Game struct {
	Version int  `json:"version"`
	Done    bool `json:"done"`

	Code        string    `json:"code"`
	Players     []*Player `json:"players"`
	Turn        int       `json:"turn"`
	Rounds      []*Round  `json:"rounds"`
	Bag         []*Tile   `json:"bag"`
	BoardWidth  int       `json:"board_width"`
	BoardHeight int       `json:"board_height"`
	MaxPips     int       `json:"max_pips"`
	History     []string  `json:"history"`
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

func (g *Game) LeaveOrQuit(name string) bool {
	if len(g.Rounds) > 0 {
		for _, p := range g.Players {
			if p.Name == name {
				g.Done = true
				return true
			}
		}
		return false
	}

	newPlayers := []*Player{}
	quitting := false
	for _, p := range g.Players {
		if p.Name == name {
			quitting = true
			break
		}
		newPlayers = append(newPlayers, p)
	}
	g.Players = newPlayers
	if len(g.Players) == 0 {
		g.Done = true
	}
	return quitting
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

func (g *Game) LastRoundLeader() int {
	if len(g.Rounds) == 0 {
		return g.MaxPips + 1
	}
	lastRound := g.Rounds[len(g.Rounds)-1]
	firstTile := lastRound.LaidTiles[0]
	return firstTile.Tile.PipsA
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

	lastRoundLeader := g.LastRoundLeader()

	playerLines := map[string][]*LaidTile{}
	for _, p := range g.Players {
		playerLines[p.Name] = []*LaidTile{}
		p.Dead = false
		p.ChickenFoot = false
		p.JustDrew = false
	}

	g.Rounds = append(g.Rounds, &Round{
		Turn:        0,
		LaidTiles:   []*LaidTile{},
		PlayerLines: playerLines,
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

	// Give each player 7 tiles.
	for _, p := range g.Players {
		p.Hand = g.Bag[:7]
		g.Bag = g.Bag[7:]
	}
	g.Turn = 0

	// Find the round leader, drawing tiles if we need to.
	foundLeader := false
	var potentialLeader int
	for !foundLeader {
		for potentialLeader = lastRoundLeader - 1; potentialLeader > 0; potentialLeader-- {
			for i, p := range g.Players {
				if !p.HasRoundLeader(potentialLeader) {
					continue
				}
				log.Printf("%s is the round leader", p.Name)
				foundLeader = true
				g.Turn = i
				break
			}
			if foundLeader {
				break
			}
		}
		if foundLeader {
			break
		}
		for _, p := range g.Players {
			g.DrawTile(p.Name)
		}
	}

	g.Note(fmt.Sprintf("%s started round %d - %d:%d", g.Players[g.Turn].Name, len(g.Rounds), potentialLeader, potentialLeader))

	if err := g.LayTile(&LaidTile{
		Tile:        &Tile{PipsA: potentialLeader, PipsB: potentialLeader},
		PlayerName:  g.Players[g.Turn].Name,
		Orientation: "right",
		X:           g.BoardWidth/2 - 1,
		Y:           g.BoardHeight / 2,
	}); err != nil {
		return fmt.Errorf("laying round leader tile: %w", err)
	}

	return nil
}

func (g *Game) Pass(name string) bool {
	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	if player == nil {
		return false
	}

	if !player.JustDrew {
		return false
	}
	g.Turn = (g.Turn + 1) % len(g.Players)
	player.JustDrew = false
	g.CurrentRound().Note(fmt.Sprintf("%s passed", name))
	return true
}

func (g *Game) Note(n string) {
	g.History = append(g.History, n)
}

func (g *Game) CurrentRound() *Round {
	if len(g.Rounds) == 0 {
		return nil
	}
	r := g.Rounds[len(g.Rounds)-1]
	if r.Done {
		return nil
	}
	return r
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
	g.Bag = g.Bag[1:]

	player.JustDrew = true

	if !player.ChickenFoot {
		player.ChickenFoot = true
		g.Note(fmt.Sprintf("%s is chicken-footed", name))
	}

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

	round := g.CurrentRound()
	if round == nil {
		return ErrRoundNotStarted
	}
	if err := round.LayTile(g, tile); err != nil {
		return fmt.Errorf("laying tile: %w", err)
	}
	g.Turn = (g.Turn + 1) % len(g.Players)

	round.Note(fmt.Sprintf("%s laid %d:%d", tile.PlayerName, tile.Tile.PipsA, tile.Tile.PipsB))
	if player.ChickenFoot {
		g.Note(fmt.Sprintf("%s is no longer chicken-footed", tile.PlayerName))
		player.ChickenFoot = false
	}

	livingPlayers := []*Player{}
	for _, p := range g.Players {
		if !p.Dead {
			livingPlayers = append(livingPlayers, p)
		}
	}
	if len(livingPlayers) == 1 {
		round.Done = true
		g.Note(fmt.Sprintf("%s wins the round", livingPlayers[0].Name))
		livingPlayers[0].Score += 2
	} else {
		for _, p := range livingPlayers {
			if len(p.Hand) == 0 {
				round.Done = true
				g.Note(fmt.Sprintf("%s wins the round", p.Name))
				p.Score += 2
			}
		}
	}

	player.JustDrew = false
	return nil
}

type Player struct {
	Name        string  `json:"name"`
	Score       int     `json:"score"`
	Hand        []*Tile `json:"hand"`
	ChickenFoot bool    `json:"chicken_foot"`
	Dead        bool    `json:"dead"`
	JustDrew    bool    `json:"just_drew"`
}

func (p *Player) HasRoundLeader(leader int) bool {
	for _, t := range p.Hand {
		if t.PipsA != t.PipsB {
			continue
		}
		if t.PipsA == leader {
			return true
		}
	}
	return false
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
	NextPips    int    `json:"next_pips"`
}

func (lt *LaidTile) CoordA() string {
	return fmt.Sprintf("%d,%d", lt.X, lt.Y)
}

func (lt *LaidTile) CoordAX() int {
	return lt.X
}

func (lt *LaidTile) CoordAY() int {
	return lt.Y
}

func (lt *LaidTile) CoordBX() int {
	switch lt.Orientation {
	case "right":
		return lt.X + 1
	case "left":
		return lt.X - 1
	}
	return lt.X
}

func (lt *LaidTile) CoordBY() int {
	switch lt.Orientation {
	case "up":
		return lt.Y - 1
	case "down":
		return lt.Y + 1
	}
	return lt.Y
}

func (lt *LaidTile) CoordB() string {
	return fmt.Sprintf("%d,%d", lt.CoordBX(), lt.CoordBY())
}

func (lt *LaidTile) String() string {
	return fmt.Sprintf("%d:%d %s:%s", lt.Tile.PipsA, lt.Tile.PipsB, lt.CoordA(), lt.CoordB())
}

type Round struct {
	Turn        int                    `json:"turn"`
	LaidTiles   []*LaidTile            `json:"laid_tiles"`
	Done        bool                   `json:"done"`
	History     []string               `json:"history"`
	PlayerLines map[string][]*LaidTile `json:"player_lines"`
	FreeLines   [][]*LaidTile          `json:"free_lines"`
}

func (r *Round) Note(n string) {
	if r == nil {
		return
	}
	r.History = append(r.History, n)
}

func (r *Round) LayTile(g *Game, lt *LaidTile) error {
	if r.Done {
		return ErrRoundAlreadyDone
	}

	if len(r.LaidTiles) == 0 {
		if lt.Tile.PipsA != lt.Tile.PipsB {
			return ErrNotRoundLeader
		}
		lt.NextPips = lt.Tile.PipsA
		lt.PlayerName = ""
		r.LaidTiles = append(r.LaidTiles, lt)
		for n, p := range r.PlayerLines {
			r.PlayerLines[n] = append(p, lt)
		}
		return nil
	}

	if lt.CoordAX() < 0 || lt.CoordAX() >= g.BoardWidth {
		return ErrTileOutOfBounds
	}
	if lt.CoordAY() < 0 || lt.CoordAY() >= g.BoardHeight {
		return ErrTileOutOfBounds
	}
	if lt.CoordBX() < 0 || lt.CoordBX() >= g.BoardWidth {
		return ErrTileOutOfBounds
	}
	if lt.CoordBY() < 0 || lt.CoordBY() >= g.BoardHeight {
		return ErrTileOutOfBounds
	}

	squarePips := r.MapTiles()
	if ot, ok := squarePips[lt.CoordA()]; ok {
		log.Printf("%s A-blocked by %s", lt, ot.LaidTile)
		return ErrTileOccluded
	}
	if ot, ok := squarePips[lt.CoordB()]; ok {
		log.Printf("%s B-blocked by %s", lt, ot.LaidTile)
		return ErrTileOccluded
	}

	canPlayOnLine := func(line []*LaidTile) (bool, int) {
		last := line[len(line)-1]
		if lt.Tile.PipsA == last.NextPips {
			if last.Tile.PipsA == lt.Tile.PipsA {
				if last.CoordAX() == lt.CoordAX() &&
					(last.CoordAY() == lt.CoordAY()+1 || last.CoordAY() == lt.CoordAY()-1) {
					return true, lt.Tile.PipsB
				}
				if last.CoordAY() == lt.CoordAY() &&
					(last.CoordAX() == lt.CoordAX()+1 || last.CoordAX() == lt.CoordAX()-1) {
					return true, lt.Tile.PipsB
				}
			}
			if last.Tile.PipsB == lt.Tile.PipsA {
				if last.CoordBX() == lt.CoordAX() &&
					(last.CoordBY() == lt.CoordAY()+1 || last.CoordBY() == lt.CoordAY()-1) {
					return true, lt.Tile.PipsB
				}
				if last.CoordBY() == lt.CoordAY() &&
					(last.CoordBX() == lt.CoordAX()+1 || last.CoordBX() == lt.CoordAX()-1) {
					return true, lt.Tile.PipsB
				}
			}
		}
		if lt.Tile.PipsB == last.NextPips {
			if last.Tile.PipsA == lt.Tile.PipsB {
				if last.CoordAX() == lt.CoordBX() &&
					(last.CoordAY() == lt.CoordBY()+1 || last.CoordAY() == lt.CoordBY()-1) {
					return true, lt.Tile.PipsA
				}
				if last.CoordAY() == lt.CoordBY() &&
					(last.CoordAX() == lt.CoordBX()+1 || last.CoordAX() == lt.CoordBX()-1) {
					return true, lt.Tile.PipsA
				}
			}
			if last.Tile.PipsB == lt.Tile.PipsB {
				if last.CoordBX() == lt.CoordBX() &&
					(last.CoordBY() == lt.CoordBY()+1 || last.CoordBY() == lt.CoordBY()-1) {
					return true, lt.Tile.PipsA
				}
				if last.CoordBY() == lt.CoordBY() &&
					(last.CoordBX() == lt.CoordBX()+1 || last.CoordBX() == lt.CoordBX()-1) {
					return true, lt.Tile.PipsA
				}
			}
		}
		return false, 0
	}

	numLinesPlayed := 0

	player := g.GetPlayer(lt.PlayerName)
	if !player.Dead {
		mainLine := r.PlayerLines[lt.PlayerName]
		if ok, nextPips := canPlayOnLine(mainLine); ok {
			numLinesPlayed++
			r.PlayerLines[lt.PlayerName] = append(mainLine, lt)
			lt.NextPips = nextPips
		}
	}

	if player.Dead || !player.ChickenFoot {
		for name, line := range r.PlayerLines {
			if name == player.Name {
				continue
			}
			if !player.ChickenFoot {
				continue
			}
			if player.Dead {
				continue
			}
			if ok, nextPips := canPlayOnLine(line); ok {
				numLinesPlayed++
				r.PlayerLines[name] = append(line, lt)
				lt.NextPips = nextPips
			}
		}
		for i, line := range r.FreeLines {
			if ok, nextPips := canPlayOnLine(line); ok {
				numLinesPlayed++
				r.FreeLines[i] = append(line, lt)
				lt.NextPips = nextPips
			}
		}
	}

	if numLinesPlayed == 0 {
		return ErrNoLine
	}

	// Add the new tile to the occlusion grid.
	squarePips[lt.CoordA()] = SquarePips{LaidTile: lt, Pips: lt.Tile.PipsA}
	squarePips[lt.CoordB()] = SquarePips{LaidTile: lt, Pips: lt.Tile.PipsB}
	isOpenFrom := func(x, y int) bool {
		c := func(x, y int) string {
			return fmt.Sprintf("%d,%d", x, y)
		}
		adj := func(x, y int) []string {
			return []string{c(x-1, y), c(x+1, y), c(x, y-1), c(x, y+1)}
		}
		if _, ok := squarePips[c(x, y+1)]; !ok {
			for _, n := range adj(x, y+1) {
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		if _, ok := squarePips[c(x, y-1)]; !ok {
			for _, n := range adj(x, y-1) {
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		if _, ok := squarePips[c(x+1, y)]; !ok {
			for _, n := range adj(x+1, y) {
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		if _, ok := squarePips[c(x-1, y)]; !ok {
			for _, n := range adj(x-1, y) {
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		return false
	}
	isCutOff := func(line []*LaidTile) bool {
		last := line[len(line)-1]
		if last.NextPips == last.Tile.PipsA && isOpenFrom(last.CoordAX(), last.CoordAY()) {
			return false
		}
		if last.NextPips == last.Tile.PipsB && isOpenFrom(last.CoordBX(), last.CoordBY()) {
			return false
		}
		return true
	}

	// Kill lines as needed
	newFreeLines := [][]*LaidTile{}
	for _, line := range r.FreeLines {
		if isCutOff(line) {
			g.Note(fmt.Sprintf("%s cut-off a free line", lt.PlayerName))
			continue
		}
		newFreeLines = append(newFreeLines, line)
	}
	r.FreeLines = newFreeLines

	for _, p := range g.Players {
		if !isCutOff(r.PlayerLines[p.Name]) {
			continue
		}
		if p.Name == player.Name {
			g.Note(fmt.Sprintf("%s cut-off their own line", player.Name))
		} else {
			g.Note(fmt.Sprintf("%s cut-off %s's line", player.Name, p.Name))
		}
		player.Score += 1
		p.Score -= 1
		p.Dead = true
		p.ChickenFoot = false
	}

	r.LaidTiles = append(r.LaidTiles, lt)
	return nil
}

type SquarePips struct {
	LaidTile *LaidTile
	Pips     int
}

func (r *Round) MapTiles() map[string]SquarePips {
	squarePips := map[string]SquarePips{}
	for _, lt := range r.LaidTiles {
		squarePips[lt.CoordA()] = SquarePips{
			LaidTile: lt,
			Pips:     lt.Tile.PipsA,
		}
		squarePips[lt.CoordB()] = SquarePips{
			LaidTile: lt,
			Pips:     lt.Tile.PipsB,
		}
	}
	return squarePips
}
