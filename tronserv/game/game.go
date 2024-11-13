package game

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

var Colors = []string{"red", "blue", "green"}

type Game struct {
	Version int64 `json:"version"`
	Done    bool  `json:"done"`

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

	switch len(g.Players) {
	case 1:
		g.BoardWidth = 6
		g.BoardHeight = 7
		g.MaxPips = 6
	case 2:
		g.BoardWidth = 8
		g.BoardHeight = 9
		g.MaxPips = 7
	case 3:
		g.BoardWidth = 10
		g.BoardHeight = 11
		g.MaxPips = 8
	case 4:
		g.BoardWidth = 12
		g.BoardHeight = 13
		g.MaxPips = 10
	case 5:
		g.BoardWidth = 14
		g.BoardHeight = 15
		g.MaxPips = 11
	case 6:
		g.BoardWidth = 16
		g.BoardHeight = 17
		g.MaxPips = 12
	}

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
	if len(g.Players) < 1 {
		return ErrGameNotEnoughPlayers
	}

	if len(g.Rounds) > 0 {
		if !g.Rounds[len(g.Rounds)-1].Done {
			return ErrGamePreviousRoundNotDone
		}
	}

	lastRoundLeader := g.LastRoundLeader()
	if lastRoundLeader == 0 {
		return ErrGameOver
	}

	playerLines := map[string][]*LaidTile{}
	for _, p := range g.Players {
		playerLines[p.Name] = []*LaidTile{}
		p.Dead = false
		p.ChickenFoot = false
		p.JustDrew = false
		p.Hand = nil
	}

	g.Rounds = append(g.Rounds, &Round{
		Turn:        0,
		LaidTiles:   []*LaidTile{},
		PlayerLines: playerLines,
	})

	// Fill the bag with tiles.z
	g.Bag = nil
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

	switch g.Code[:6] {
	case "AAAAAA":
		log.Print("secret code")
		g.Bag = []*Tile{{
			PipsA: 1, PipsB: 1,
		}, {
			PipsA: 1, PipsB: 2,
		}, {
			PipsA: 2, PipsB: 4,
		}, {
			PipsA: 4, PipsB: 6,
		}, {
			PipsA: 6, PipsB: 8,
		}, {
			PipsA: 8, PipsB: 10,
		}, {
			PipsA: 10, PipsB: 12,
		}, {
			PipsA: 1, PipsB: 3,
		}, {
			PipsA: 3, PipsB: 5,
		}, {
			PipsA: 5, PipsB: 7,
		}, {
			PipsA: 7, PipsB: 9,
		}, {
			PipsA: 9, PipsB: 11,
		}, {
			PipsA: 11, PipsB: 13,
		}, {
			PipsA: 13, PipsB: 15,
		}, {
			PipsA: 0, PipsB: 0,
		}}
		for i := 2; i <= g.MaxPips; i++ {
			g.Bag = append(g.Bag, &Tile{PipsA: i, PipsB: i})
		}
	case "BBBBBB":
		log.Print("secret code")
		g.Bag = []*Tile{{
			PipsA: 0, PipsB: 0,
		}, {
			PipsA: 0, PipsB: 1,
		}, {
			PipsA: 0, PipsB: 2,
		}, {
			PipsA: 0, PipsB: 3,
		}, {
			PipsA: 0, PipsB: 4,
		}, {
			PipsA: 0, PipsB: 5,
		}, {
			PipsA: 0, PipsB: 6,
		}}
		for i := 1; i <= g.MaxPips; i++ {
			g.Bag = append(g.Bag, &Tile{PipsA: i, PipsB: i})
		}
	}

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
		for potentialLeader = lastRoundLeader - 1; potentialLeader >= 0; potentialLeader-- {
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

	if err := g.LayTile(g.Players[g.Turn].Name, &LaidTile{
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

func (g *Game) Pass(name string, chickenFootX, chickenFootY int) error {
	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	if player == nil {
		return ErrNoSuchPlayer
	}

	if !player.JustDrew {
		return ErrMustDrawTile
	}
	g.Turn = (g.Turn + 1) % len(g.Players)
	player.JustDrew = false
	r := g.CurrentRound()
	if r == nil {
		return ErrRoundNotStarted
	}
	r.Note(fmt.Sprintf("%s passed", name))

	if !player.ChickenFoot && !player.Dead {
		player.ChickenFoot = true
		g.Note(fmt.Sprintf("%s is chicken-footed", name))

		mainLine := r.PlayerLines[player.Name]
		if len(mainLine) > 1 {
			mostRecent := mainLine[len(mainLine)-1]
			if mostRecent.NextPips == mostRecent.Tile.PipsA {
				player.ChickenFootX = mostRecent.CoordBX()
				player.ChickenFootY = mostRecent.CoordBY()
			} else {
				player.ChickenFootX = mostRecent.CoordAX()
				player.ChickenFootY = mostRecent.CoordAY()
			}
		} else {
			// we have to pick a viable spot left around the round leader
			if chickenFootX == -1 || chickenFootY == -1 {
				return ErrMustPickChickenFoot
			}
			player.ChickenFootX = chickenFootX
			player.ChickenFootY = chickenFootY
		}
	}
	return nil
}

func (g *Game) Note(n string) {
	g.History = append(g.History, n)
	log.Print(n)
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
	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	if player == nil {
		return false
	}

	if len(g.Bag) > 0 {
		player.Hand = append(player.Hand, g.Bag[0])
		g.Bag = g.Bag[1:]
	}

	player.JustDrew = true

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

func (g *Game) LayTile(name string, tile *LaidTile) error {
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

	firstTile := len(round.LaidTiles) == 0
	if err := round.LayTile(g, name, tile, false); err != nil {
		return err
	}
	if firstTile || tile.Tile.PipsA != tile.Tile.PipsB {
		g.Turn = (g.Turn + 1) % len(g.Players)
	}

	round.Note(fmt.Sprintf("%s laid %d:%d", name, tile.Tile.PipsA, tile.Tile.PipsB))
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
	if len(livingPlayers) == 1 && len(g.Players) > 1 {
		round.Done = true
		g.Note(fmt.Sprintf("%s wins the round", livingPlayers[0].Name))
		livingPlayers[0].Score += 2
	} else if len(livingPlayers) == 0 {
		round.Done = true
		g.Note("you win I guess")
	} else {
		for _, p := range livingPlayers {
			if len(p.Hand) == 0 {
				round.Done = true
				g.Note(fmt.Sprintf("%s wins the round", p.Name))
				p.Score += 2
			}
		}
	}

	if round.Done && g.LastRoundLeader() == 0 {
		g.Done = true
	}

	player.JustDrew = false

	return nil
}

type Player struct {
	Name         string     `json:"name"`
	Score        int        `json:"score"`
	Hand         []*Tile    `json:"hand"`
	Hints        [][]string `json:"hints"`
	ChickenFoot  bool       `json:"chicken_foot"`
	Dead         bool       `json:"dead"`
	JustDrew     bool       `json:"just_drew"`
	ChickenFootX int        `json:"chicken_foot_x"`
	ChickenFootY int        `json:"chicken_foot_y"`
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

func (t *Tile) String() string {
	return fmt.Sprintf("%d:%d", t.PipsA, t.PipsB)
}

type LaidTile struct {
	Tile        *Tile  `json:"tile"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Orientation string `json:"orientation"`
	PlayerName  string `json:"player_name"`
	NextPips    int    `json:"next_pips"`
	Dead        bool   `json:"dead"`
	Indicated   *Tile  `json:"indicated"`
}

func (lt *LaidTile) Reverse() *LaidTile {
	rt := &LaidTile{}
	*rt = *lt
	rt.X, rt.Y = rt.CoordBX(), rt.CoordBY()
	switch lt.Orientation {
	case "up":
		rt.Orientation = "down"
	case "down":
		rt.Orientation = "up"
	case "left":
		rt.Orientation = "right"
	case "right":
		rt.Orientation = "left"
	}
	return rt
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

func (r *Round) FindHints(g *Game, name string, p *Player) {
	squarePips := r.MapTiles()

	hints := make([]map[string]bool, len(p.Hand))
	for i := range hints {
		hints[i] = map[string]bool{}
	}

	hintAt := func(i int, coord string) {
		hints[i][coord] = true
	}

	for i, t := range p.Hand {
		movesOffSquare := func(head *LaidTile, x, y int, op *Player) {

			for _, orientation := range []string{"up", "down", "left", "right"} {
				lt := &LaidTile{
					Tile:        t,
					Orientation: orientation,
					X:           x,
					Y:           y,
				}
				if r.LayTile(g, name, lt, true) == nil || r.LayTile(g, name, lt.Reverse(), true) == nil {
					hintAt(i, lt.CoordA())
					hintAt(i, lt.CoordB())
				}
			}

		}

		movesOffTile := func(head *LaidTile, op *Player) {
			movesOffSquare(head, head.CoordAX()-1, head.CoordAY(), op)
			movesOffSquare(head, head.CoordAX()+1, head.CoordAY(), op)
			movesOffSquare(head, head.CoordAX(), head.CoordAY()-1, op)
			movesOffSquare(head, head.CoordAX(), head.CoordAY()+1, op)
			movesOffSquare(head, head.CoordBX()-1, head.CoordBY(), op)
			movesOffSquare(head, head.CoordBX()+1, head.CoordBY(), op)
			movesOffSquare(head, head.CoordBX(), head.CoordBY()-1, op)
			movesOffSquare(head, head.CoordBX(), head.CoordBY()+1, op)
		}
		p := g.GetPlayer(name)

		// first consider direct plays
		for opname, line := range r.PlayerLines {
			op := g.GetPlayer(opname)
			if opname != name {
				if p.ChickenFoot || p.Dead {
					continue
				}
				if !op.ChickenFoot {
					continue
				}
			}
			movesOffTile(line[len(line)-1], op)
		}
		if p.ChickenFoot {
			// no free line activity allowed
			return
		}

		// then consider new free lines
		if t.PipsA != t.PipsB {
			continue
		}
		if t.PipsA < r.PlayerLines[name][0].Tile.PipsA {
			continue
		}
		isHigher := true
		for _, l := range r.FreeLines {
			if t.PipsA < l[0].Tile.PipsA {
				isHigher = false
				break
			}
		}
		if !isHigher {
			continue
		}

		log.Printf("potential free liner: %s", t)
		// potential free liner
		tryFromCoord := func(x1, y1 int) {
			log.Printf(" trying from %d,%d", x1, y1)
			if x1 < 0 || x1 >= g.BoardWidth {
				return
			}
			if y1 < 0 || y1 >= g.BoardHeight {
				return
			}
			if _, ok := squarePips[fmt.Sprintf("%d,%d", x1, y1)]; ok {
				return
			}
			tryToCoord := func(x2, y2 int) {
				log.Printf(" trying to %d,%d", x2, y2)
				if x2 < 0 || x2 >= g.BoardWidth {
					return
				}
				if y2 < 0 || y2 >= g.BoardHeight {
					return
				}
				if x2 < x1 {
					x1, x2 = x2, x1
				}
				if y2 < y1 {
					y1, y2 = y2, y1
				}
				if !sixPathFrom(squarePips, x1, y1, x2, y2) {
					return
				}
				tryA := func(x, y int) {
					for _, orientation := range []string{"up", "down", "left", "right"} {
						lt := &LaidTile{
							Tile:        t,
							Orientation: orientation,
							X:           x,
							Y:           y,
						}
						if r.LayTile(g, name, lt, true) == nil || r.LayTile(g, name, lt.Reverse(), true) == nil {
							hintAt(i, lt.CoordA())
							hintAt(i, lt.CoordB())
						}
					}
				}
				tryA(x2+1, y2)
				tryA(x2-1, y2)
				tryA(x2, y2+1)
				tryA(x2, y2-1)
			}
			tryToCoord(x1+5, y1)
			tryToCoord(x1-5, y1)
			tryToCoord(x1, y1+5)
			tryToCoord(x1, y1-5)
		}
		tryFreeFrom := func(head *LaidTile) {
			if head.NextPips == t.PipsA {
				tryFromCoord(head.CoordAX()-1, head.CoordAY())
				tryFromCoord(head.CoordAX()+1, head.CoordAY())
				tryFromCoord(head.CoordAX(), head.CoordAY()-1)
				tryFromCoord(head.CoordAX(), head.CoordAY()+1)
			}
			if head.NextPips == t.PipsB {
				tryFromCoord(head.CoordBX()-1, head.CoordBY())
				tryFromCoord(head.CoordBX()+1, head.CoordBY())
				tryFromCoord(head.CoordBX(), head.CoordBY()-1)
				tryFromCoord(head.CoordBX(), head.CoordBY()+1)
			}
		}
		for _, l := range r.PlayerLines {
			tryFreeFrom(l[len(l)-1])
		}
		for _, l := range r.FreeLines {
			tryFreeFrom(l[len(l)-1])
		}
	}
	p.Hints = make([][]string, len(p.Hand))
	for i, hintList := range hints {
		for h := range hintList {
			p.Hints[i] = append(p.Hints[i], h)
		}
	}
}

func (r *Round) canPlayOnLine(lt *LaidTile, line []*LaidTile) (bool, int) {
	last := line[len(line)-1]
	return r.canPlayOnTile(lt, last)
}

func (r *Round) canPlayOnTile(lt, last *LaidTile) (bool, int) {
	if lt.Indicated != nil && lt.Indicated.PipsA != -1 {
		if last.Tile.PipsA != lt.Indicated.PipsA || last.Tile.PipsB != lt.Indicated.PipsB {
			return false, 0
		}
	}
	return r.canPlayOnTileWithoutIndication(lt, last)
}

func (r *Round) canPlayOnTileWithoutIndication(lt, last *LaidTile) (bool, int) {
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

func sixPathFrom(squarePips map[string]SquarePips, x1, y1, x2, y2 int) bool {
	switch {
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		if y2-y1 != 5 {
			return false
		}
		for y := y1; y <= y2; y++ {
			if _, ok := squarePips[fmt.Sprintf("%d,%d", x1, y)]; ok {
				return false
			}
		}
		return true
	case y1 == y2:
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if x2-x1 != 5 {
			return false
		}
		for x := x1; x <= x2; x++ {
			if _, ok := squarePips[fmt.Sprintf("%d,%d", x, y1)]; ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (r *Round) LayTile(g *Game, name string, lt *LaidTile, dryRun bool) error {
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
	if _, ok := squarePips[lt.CoordA()]; ok {
		return ErrTileOccluded
	}
	if _, ok := squarePips[lt.CoordB()]; ok {
		return ErrTileOccluded
	}

	isInRoundLeaderChickenfoot := false
	for _, p := range g.Players {
		if !p.ChickenFoot {
			continue
		}
		if p.Name == name {
			// It's ok to play on our own foot.
			continue
		}
		// this only matters for that initial chickenfoot.
		if len(r.PlayerLines[p.Name]) > 1 {
			continue
		}
		// we don't check if this is a round-leader chickenfoot, because
		// otherwise it would be blocked by a tile and already illegal.
		if lt.CoordAX() == p.ChickenFootX && lt.CoordAY() == p.ChickenFootY {
			isInRoundLeaderChickenfoot = true
			break
		}
		if lt.CoordBX() == p.ChickenFootX && lt.CoordBY() == p.ChickenFootY {
			isInRoundLeaderChickenfoot = true
			break
		}
		// we're not playing on this chicken-foot, so we can't block it.
		available := map[string]bool{}
		checkCanPlayOn := func(x, y int) bool {
			if x < 0 || x >= g.BoardWidth || y < 0 || y >= g.BoardHeight {
				return false
			}
			key := fmt.Sprintf("%d,%d", x, y)
			if _, ok := squarePips[key]; ok {
				return false
			}
			available[key] = true
			return true
		}
		checkCanPlayOn(p.ChickenFootX+1, p.ChickenFootY)
		checkCanPlayOn(p.ChickenFootX-1, p.ChickenFootY)
		checkCanPlayOn(p.ChickenFootX, p.ChickenFootY+1)
		checkCanPlayOn(p.ChickenFootX, p.ChickenFootY-1)
		if len(available) == 1 {
			for k := range available {
				if lt.CoordA() == k || lt.CoordB() == k {
					return ErrNoBlockingFeet
				}
			}
		}
	}

	playedALine := false

	player := g.GetPlayer(name)
	if !player.Dead && !isInRoundLeaderChickenfoot {
		mainLine := r.PlayerLines[player.Name]

		onFoot := false
		if player.ChickenFoot && len(mainLine) == 1 {
			if player.ChickenFootX == lt.CoordAX() && player.ChickenFootY == lt.CoordAY() {
				onFoot = true
			}
			if player.ChickenFootX == lt.CoordBX() && player.ChickenFootY == lt.CoordBY() {
				onFoot = true
			}
		}
		if len(mainLine) > 1 || onFoot || !player.ChickenFoot {
			if ok, nextPips := r.canPlayOnLine(lt, mainLine); ok {
				playedALine = true
				if !dryRun {
					r.PlayerLines[player.Name] = append(mainLine, lt)
					lt.NextPips = nextPips
				}
			}
		}
	}
	if !player.Dead || !player.ChickenFoot {
		for oname, line := range r.PlayerLines {
			if playedALine {
				continue
			}
			op := g.GetPlayer(oname)
			if oname == player.Name {
				continue
			}
			if !op.ChickenFoot {
				continue
			}
			if op.Dead {
				continue
			}
			if len(line) == 1 {
				// round leader, need to play on top of chickenfoot.
				isOnMyFoot := false
				if lt.CoordAX() == op.ChickenFootX && lt.CoordAY() == op.ChickenFootY && lt.Tile.PipsA == line[0].NextPips {
					isOnMyFoot = true
				}
				if lt.CoordBX() == op.ChickenFootX && lt.CoordBY() == op.ChickenFootY && lt.Tile.PipsB == line[0].NextPips {
					isOnMyFoot = true
				}
				if !isOnMyFoot {
					continue
				}
			}
			if ok, nextPips := r.canPlayOnLine(lt, line); ok {
				playedALine = true
				if !dryRun {
					lt.PlayerName = oname
					r.PlayerLines[oname] = append(line, lt)
					lt.NextPips = nextPips
					// Update the chicken-foot.
					if lt.NextPips == lt.Tile.PipsA {
						op.ChickenFootX = lt.CoordBX()
						op.ChickenFootY = lt.CoordBY()
					}
					if lt.NextPips == lt.Tile.PipsB {
						op.ChickenFootX = lt.CoordAX()
						op.ChickenFootY = lt.CoordAY()
					}
				}
			}
		}
		for i, line := range r.FreeLines {
			if playedALine {
				continue
			}
			if ok, nextPips := r.canPlayOnLine(lt, line); ok {
				playedALine = true
				if !dryRun {
					lt.PlayerName = ""
					r.FreeLines[i] = append(line, lt)
					lt.NextPips = nextPips
				}
			}
		}
	}

	if (player.Dead || !player.ChickenFoot) && lt.Tile.PipsA == lt.Tile.PipsB {
		isHigher := true
		if lt.Tile.PipsA < r.PlayerLines[g.Players[0].Name][0].Tile.PipsA {
			isHigher = false
		}
		for _, l := range r.FreeLines {
			if lt.Tile.PipsA < l[0].Tile.PipsA {
				isHigher = false
			}
		}
		// can we start a free line
		canStartFreeLineOff := func(head *LaidTile) bool {
			type pair struct {
				x, y int
			}
			pairsHead := []pair{{
				head.CoordAX() - 1, head.CoordAY(),
			}, {
				head.CoordAX() + 1, head.CoordAY(),
			}, {
				head.CoordAX(), head.CoordAY() - 1,
			}, {
				head.CoordAX(), head.CoordAY() + 1,
			}, {
				head.CoordBX() - 1, head.CoordBY(),
			}, {
				head.CoordBX() + 1, head.CoordBY(),
			}, {
				head.CoordBX(), head.CoordBY() - 1,
			}, {
				head.CoordBX(), head.CoordBY() + 1,
			}}
			pairsLT := []pair{{
				lt.CoordAX() - 1, lt.CoordAY(),
			}, {
				lt.CoordAX() + 1, lt.CoordAY(),
			}, {
				lt.CoordAX(), lt.CoordAY() - 1,
			}, {
				lt.CoordAX(), lt.CoordAY() + 1,
			}, {
				lt.CoordBX() - 1, lt.CoordBY(),
			}, {
				lt.CoordBX() + 1, lt.CoordBY(),
			}, {
				lt.CoordBX(), lt.CoordBY() - 1,
			}, {
				lt.CoordBX(), lt.CoordBY() + 1,
			}}
			for _, headPair := range pairsHead {
				for _, ltPair := range pairsLT {
					if lt.CoordAX() >= headPair.x && lt.CoordAX() <= ltPair.x && lt.CoordAY() >= headPair.y && lt.CoordAY() <= ltPair.y {
						continue
					}
					if lt.CoordBX() >= headPair.x && lt.CoordBX() <= ltPair.x && lt.CoordBY() >= headPair.y && lt.CoordBY() <= ltPair.y {
						continue
					}
					if sixPathFrom(squarePips, headPair.x, headPair.y, ltPair.x, ltPair.y) {
						return true
					}
				}
			}
			return false
		}
		if isHigher {
			canBeFree := false
			for _, l := range r.PlayerLines {
				if canStartFreeLineOff(l[len(l)-1]) {
					canBeFree = true
					break
				}
			}
			for _, l := range r.FreeLines {
				if canStartFreeLineOff(l[len(l)-1]) {
					canBeFree = true
					break
				}
			}
			if canBeFree {
				playedALine = true
				if !dryRun {
					lt.PlayerName = ""
					r.FreeLines = append(r.FreeLines, []*LaidTile{lt})
					lt.NextPips = lt.Tile.PipsA
					g.Note(fmt.Sprintf("%s started a free line", name))
				}
			}
		}
	}

	if !playedALine {
		return ErrNoLine
	}

	// Add the new tile to the occlusion grid.
	squarePips[lt.CoordA()] = SquarePips{LaidTile: lt, Pips: lt.Tile.PipsA}
	squarePips[lt.CoordB()] = SquarePips{LaidTile: lt, Pips: lt.Tile.PipsB}
	isOpenFrom := func(x, y int) bool {
		c := func(x, y int) string {
			if x < 0 || x >= g.BoardWidth || y < 0 || y >= g.BoardHeight {
				return ""
			}
			return fmt.Sprintf("%d,%d", x, y)
		}
		adj := func(x, y int) []string {
			return []string{c(x-1, y), c(x+1, y), c(x, y-1), c(x, y+1)}
		}
		if _, ok := squarePips[c(x, y+1)]; !ok {
			for _, n := range adj(x, y+1) {
				if n == "" {
					continue
				}
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		if _, ok := squarePips[c(x, y-1)]; !ok {
			for _, n := range adj(x, y-1) {
				if n == "" {
					continue
				}
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		if _, ok := squarePips[c(x+1, y)]; !ok {
			for _, n := range adj(x+1, y) {
				if n == "" {
					continue
				}
				if _, ok := squarePips[n]; !ok {
					return true
				}
			}
		}
		if _, ok := squarePips[c(x-1, y)]; !ok {
			for _, n := range adj(x-1, y) {
				if n == "" {
					continue
				}
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

	deadTileKeys := map[string]bool{}

	if !dryRun {
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
			if p.Dead {
				continue
			}
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
			for _, lt := range r.PlayerLines[p.Name][1:] {
				lt.Dead = true
				deadTileKeys[lt.Tile.String()] = true
			}
			p.ChickenFoot = false
		}

		r.LaidTiles = append(r.LaidTiles, lt)
	}

	// look for il ouroboros
	if !dryRun && !player.ChickenFoot {
		consumed := []string{}

		canConsume := func(head *LaidTile) bool {
			if lt.NextPips != head.NextPips {
				return false
			}
			checkLTCoord := func(x, y int) bool {
				adjacent := func(x2, y2 int) bool {
					if x == x2 && (y == y2-1 || y == y2+1) {
						return true
					}
					if y == y2 && (x == x2-1 || x == x2+1) {
						return true
					}
					return false
				}
				if head.Tile.PipsA == head.NextPips {
					if adjacent(head.CoordAX(), head.CoordAY()) {
						return true
					}
				}
				if head.Tile.PipsB == head.NextPips {
					if adjacent(head.CoordBX(), head.CoordBY()) {
						return true
					}
				}
				return false
			}
			if lt.Tile.PipsA == lt.NextPips {
				if checkLTCoord(lt.CoordAX(), lt.CoordAY()) {
					return true
				}
			}
			if lt.Tile.PipsB == lt.NextPips {
				if checkLTCoord(lt.CoordBX(), lt.CoordBY()) {
					return true
				}
			}
			return false
		}

		for _, op := range g.Players {
			if !op.ChickenFoot && op.Name != name {
				continue
			}
			if op.Dead {
				continue
			}
			playerLine := r.PlayerLines[op.Name]
			head := playerLine[len(playerLine)-1]
			if lt == head {
				continue
			}

			if canConsume(head) {
				consumed = append(consumed, op.Name)
			}
		}

		freeLines := 0
		newFreeLines := [][]*LaidTile{}
		for _, fl := range r.FreeLines {
			head := fl[len(fl)-1]
			if head != lt && canConsume(head) {
				freeLines += 1
				for _, lt := range fl {
					lt.Dead = true
					deadTileKeys[lt.Tile.String()] = true
				}
			} else {
				newFreeLines = append(newFreeLines, fl)
			}
		}
		r.FreeLines = newFreeLines

		if len(consumed) > 0 {
			if lt.PlayerName != "" {
				consumed = append(consumed, lt.PlayerName)
			}
			freeLineNote := ""
			if freeLines > 0 {
				freeLineNote = fmt.Sprintf("%d free line(s)", freeLines)
			}
			g.Note(fmt.Sprintf("%s's IL OUROBOROS consumes %s", name, strings.Join(append(consumed, freeLineNote), ", ")))
			for _, n := range consumed {
				op := g.GetPlayer(n)
				op.Dead = true
				op.ChickenFoot = false
				for _, lt := range r.PlayerLines[op.Name][1:] {
					lt.Dead = true
					deadTileKeys[lt.Tile.String()] = true
				}
				op.Score -= 1
				player.Score += 1
			}
		}
	}

	for _, lt := range r.LaidTiles {
		if deadTileKeys[lt.Tile.String()] {
			lt.Dead = true
		}
	}

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
