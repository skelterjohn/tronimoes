package game

import (
	"fmt"
	"log"

	"math/rand"
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
		g.MaxPips = 8
	case 3:
		g.BoardWidth = 10
		g.BoardHeight = 11
		g.MaxPips = 10
	case 4:
		g.BoardWidth = 12
		g.BoardHeight = 13
		g.MaxPips = 12
	case 5:
		g.BoardWidth = 14
		g.BoardHeight = 15
		g.MaxPips = 14
	case 6:
		g.BoardWidth = 16
		g.BoardHeight = 17
		g.MaxPips = 16
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
	r := g.CurrentRound()
	if r == nil {
		return false
	}
	r.Note(fmt.Sprintf("%s passed", name))

	if !player.ChickenFoot {
		player.ChickenFoot = true
		g.Note(fmt.Sprintf("%s is chicken-footed", name))

		mainLine, ok := r.PlayerLines[player.Name]
		if !ok {
			return false
		}
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
			roundLeader := mainLine[0]
			squarePips := r.MapTiles()
			checkFoot := func(x, y int) bool {
				if _, occupied := squarePips[fmt.Sprintf("%d,%d", x, y)]; occupied {
					return false
				}
				consider := func(nx, ny int) bool {
					if _, occupied := squarePips[fmt.Sprintf("%d,%d", nx, ny)]; occupied {
						return false
					}
					return true
				}
				if consider(x+1, y) || consider(x-1, y) || consider(x, y+1) || consider(x, y-1) {
					player.ChickenFootX = x
					player.ChickenFootY = y
					return true
				}

				return false
			}
			if !checkFoot(roundLeader.CoordAX()-1, roundLeader.CoordAY()) &&
				!checkFoot(roundLeader.CoordAX()+1, roundLeader.CoordAY()) &&
				!checkFoot(roundLeader.CoordAX(), roundLeader.CoordAY()+1) &&
				!checkFoot(roundLeader.CoordAX(), roundLeader.CoordAY()-1) &&
				!checkFoot(roundLeader.CoordBX()-1, roundLeader.CoordBY()) &&
				!checkFoot(roundLeader.CoordBX()+1, roundLeader.CoordBY()) &&
				!checkFoot(roundLeader.CoordBX(), roundLeader.CoordBY()+1) &&
				!checkFoot(roundLeader.CoordBX(), roundLeader.CoordBY()-1) {
				log.Printf("In game %q, unable to find a round-leader chicken foot for %s", g.Code, player.Name)
				return false
			}

		}
	}
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
	if err := round.LayTile(g, name, tile); err != nil {
		return fmt.Errorf("laying tile: %w", err)
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

type LaidTile struct {
	Tile        *Tile  `json:"tile"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Orientation string `json:"orientation"`
	PlayerName  string `json:"player_name"`
	NextPips    int    `json:"next_pips"`
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
	hints := make([]map[string]bool, len(p.Hand))
	for i := range hints {
		hints[i] = map[string]bool{}
	}

	squarePips := r.MapTiles()

	for i, t := range p.Hand {
		movesOffSquare := func(head *LaidTile, x, y int) {
			for _, orientation := range []string{"up", "down", "left", "right"} {
				lt := &LaidTile{
					Tile:        t,
					Orientation: orientation,
					X:           x,
					Y:           y,
				}
				// check if this tile is blocked
				if _, ok := squarePips[lt.CoordA()]; ok {
					continue
				}
				if lt.CoordAX() < 0 || lt.CoordAX() >= g.BoardWidth {
					continue
				}
				if lt.CoordAY() < 0 || lt.CoordAY() >= g.BoardHeight {
					continue
				}
				if _, ok := squarePips[lt.CoordB()]; ok {
					continue
				}
				if lt.CoordBX() < 0 || lt.CoordBX() >= g.BoardWidth {
					continue
				}
				if lt.CoordBY() < 0 || lt.CoordBY() >= g.BoardHeight {
					continue
				}
				if ok, _ := r.canPlayOnTile(lt, head); ok {
					hints[i][lt.CoordA()] = true
					hints[i][lt.CoordB()] = true
				}
				rt := lt.Reverse()
				if ok, _ := r.canPlayOnTile(rt, head); ok {
					hints[i][rt.CoordA()] = true
					hints[i][rt.CoordB()] = true
				}
			}

		}

		movesOffTile := func(head *LaidTile) {
			movesOffSquare(head, head.CoordAX()-1, head.CoordAY())
			movesOffSquare(head, head.CoordAX()+1, head.CoordAY())
			movesOffSquare(head, head.CoordAX(), head.CoordAY()-1)
			movesOffSquare(head, head.CoordAX(), head.CoordAY()+1)
			movesOffSquare(head, head.CoordBX()-1, head.CoordBY())
			movesOffSquare(head, head.CoordBX()+1, head.CoordBY())
			movesOffSquare(head, head.CoordBX(), head.CoordBY()-1)
			movesOffSquare(head, head.CoordBX(), head.CoordBY()+1)
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
			movesOffTile(line[len(line)-1])
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

		// potential free liner

		tryFromCoord := func(x1, y1 int) {

			tryToCoord := func(x2, y2 int) {
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
				log.Printf("path from %d,%d to %d,%d", x1, y1, x2, y2)
				isOpen := func(x, y int) bool {
					if x < 0 || y < 0 || x >= g.BoardWidth || y >= g.BoardHeight {
						return false
					}
					if _, ok := squarePips[fmt.Sprintf("%d,%d", x, y)]; ok {
						return false
					}
					if x1 <= x && x <= x2 && y1 <= y && y <= y2 {
						return false
					}
					return true
				}
				tryA := func(x, y int) {
					log.Printf("trying A=%d,%d", x, y)
					if !isOpen(x, y) {
						return
					}
					mark := func(x2, y2 int) {
						log.Printf("trying B=%d,%d", x2, y2)
						if isOpen(x2, y2) {
							hints[i][fmt.Sprintf("%d,%d", x, y)] = true
							hints[i][fmt.Sprintf("%d,%d", x2, y2)] = true
						}
					}
					mark(x+1, y)
					mark(x-1, y)
					mark(x, y+1)
					mark(x, y-1)
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
			tryFromCoord(head.CoordAX()-1, head.CoordAY())
			tryFromCoord(head.CoordAX()+1, head.CoordAY())
			tryFromCoord(head.CoordAX(), head.CoordAY()-1)
			tryFromCoord(head.CoordAX(), head.CoordAY()+1)
		}
		for _, l := range r.PlayerLines {
			tryFreeFrom(l[0])
		}
		for _, l := range r.FreeLines {
			tryFreeFrom(l[0])
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

func (r *Round) LayTile(g *Game, name string, lt *LaidTile) error {
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

	isInRoundLeaderChickenfoot := false
	for _, p := range g.Players {
		if !p.ChickenFoot {
			continue
		}
		if p.Name == name {
			// It's ok to play on our own foot.
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
	}

	playedALine := false

	player := g.GetPlayer(lt.PlayerName)
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
				r.PlayerLines[player.Name] = append(mainLine, lt)
				lt.NextPips = nextPips
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
				if lt.CoordAX() == op.ChickenFootX && lt.CoordAY() == op.ChickenFootY {
					isOnMyFoot = true
				}
				if lt.CoordBX() == op.ChickenFootX && lt.CoordBY() == op.ChickenFootY {
					isOnMyFoot = true
				}
				if !isOnMyFoot {
					continue
				}
			}
			if ok, nextPips := r.canPlayOnLine(lt, line); ok {
				playedALine = true
				lt.PlayerName = oname
				r.PlayerLines[oname] = append(line, lt)
				lt.NextPips = nextPips
			}
		}
		for i, line := range r.FreeLines {
			if playedALine {
				continue
			}
			if ok, nextPips := r.canPlayOnLine(lt, line); ok {
				playedALine = true
				lt.PlayerName = ""
				r.FreeLines[i] = append(line, lt)
				lt.NextPips = nextPips
			}
		}
	}

	if lt.Tile.PipsA == lt.Tile.PipsB {
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
				lt.PlayerName = ""
				r.FreeLines = append(r.FreeLines, []*LaidTile{lt})
				lt.NextPips = lt.Tile.PipsA
				g.Note(fmt.Sprintf("%s started a free line", name))
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
