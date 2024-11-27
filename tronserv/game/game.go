package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

var Colors = []string{"red", "blue", "green"}

type Game struct {
	Created int64 `json:"created"`

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

func NewGame(ctx context.Context, code string) *Game {
	return &Game{
		Created:     time.Now().Unix(),
		Code:        code,
		Version:     0,
		BoardWidth:  10,
		BoardHeight: 11,
		MaxPips:     16,
	}
}

func (g *Game) CheckForDupes(ctx context.Context, when string) {
	if g.CurrentRound(ctx) == nil {
		return
	}
	seen := map[string]bool{}
	anyDupes := false
	visit := func(t *Tile, where string) {
		if seen[t.String()] {
			log.Printf("dupe tile: %s in %s", t.String(), where)
			anyDupes = true
		}
		seen[t.String()] = true
	}
	for _, t := range g.Bag {
		visit(t, "bag")
	}
	for _, lt := range g.CurrentRound(ctx).LaidTiles {
		visit(lt.Tile, "laid tiles")
	}
	for _, p := range g.Players {
		for _, t := range p.Hand {
			visit(t, p.Name)
		}
	}
	if anyDupes {
		data, _ := json.Marshal(g)
		log.Printf("dupes during %s: %s", when, string(data))
	}
}

func (g *Game) LeaveOrQuit(ctx context.Context, name string) bool {
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
			continue
		}
		newPlayers = append(newPlayers, p)
	}
	g.Players = newPlayers
	if len(g.Players) == 0 {
		g.Done = true
	}
	return quitting
}

func (g *Game) AddPlayer(ctx context.Context, player *Player) error {
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

func (g *Game) LastRoundLeader(ctx context.Context) int {
	if len(g.Rounds) == 0 {
		return g.MaxPips + 1
	}
	lastRound := g.Rounds[len(g.Rounds)-1]
	firstTile := lastRound.LaidTiles[0]
	return firstTile.Tile.PipsA
}

func (g *Game) Start(ctx context.Context) error {
	if len(g.Players) < 1 {
		return ErrGameNotEnoughPlayers
	}

	if len(g.Rounds) > 0 {
		if !g.Rounds[len(g.Rounds)-1].Done {
			return ErrGamePreviousRoundNotDone
		}
	}

	lastRoundLeader := g.LastRoundLeader(ctx)
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
		p.Kills = nil
	}

	g.Rounds = append(g.Rounds, &Round{
		Turn:        0,
		LaidTiles:   []*LaidTile{},
		PlayerLines: playerLines,
	})

	// Fill the bag with tiles
	g.Bag = nil
	for a := 0; a <= g.MaxPips; a++ {
		for b := a; b <= g.MaxPips; b++ {
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
		}, {
			PipsA: 2, PipsB: 3,
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
			if len(g.Bag) == 0 {
				return ErrEmptyBag
			}
			p.Hand = append(p.Hand, g.Bag[0])
			g.Bag = g.Bag[1:]
		}
	}

	g.Note(ctx, fmt.Sprintf("%s started round %d - %d:%d", g.Players[g.Turn].Name, len(g.Rounds), potentialLeader, potentialLeader))

	if err := g.LayTile(ctx, g.Players[g.Turn].Name, &LaidTile{
		Tile:        &Tile{PipsA: potentialLeader, PipsB: potentialLeader},
		PlayerName:  g.Players[g.Turn].Name,
		Orientation: "right",
		Coord: Coord{
			X: g.BoardWidth/2 - 1,
			Y: g.BoardHeight / 2,
		},
	}); err != nil {
		return fmt.Errorf("laying round leader tile: %w", err)
	}

	return nil
}

func (g *Game) Pass(ctx context.Context, name string, chickenFootX, chickenFootY int) error {
	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	if player == nil {
		return ErrNoSuchPlayer
	}

	if len(g.Bag) > 0 && !player.JustDrew {
		return ErrMustDrawTile
	}

	round := g.CurrentRound(ctx)

	if !player.JustDrew && len(g.Bag) == 0 {
		round.BaglessPasses++
	} else {
		round.BaglessPasses = 0
	}

	if round != nil {
		round.Spacer = nil
	}
	g.Turn = (g.Turn + 1) % len(g.Players)
	player.JustDrew = false
	r := g.CurrentRound(ctx)
	if r == nil {
		return ErrRoundNotStarted
	}

	chickenFootMessage := ""
	if !player.ChickenFoot && !player.Dead {
		player.ChickenFoot = true
		chickenFootMessage = " and is on the foot"

		mainLine := r.PlayerLines[player.Name]
		if len(mainLine) > 1 {
			mostRecent := mainLine[len(mainLine)-1]
			if mostRecent.NextPips == mostRecent.Tile.PipsA {
				player.ChickenFootCoord = mostRecent.CoordB()
			} else {
				player.ChickenFootCoord = mostRecent.CoordA()
			}
		} else {
			// we have to pick a viable spot left around the round leader
			if chickenFootX == -1 || chickenFootY == -1 {
				return ErrMustPickChickenFoot
			}
			player.ChickenFootCoord = Coord{X: chickenFootX, Y: chickenFootY}
		}
	}

	if round.BaglessPasses != 0 {
		r.Note(ctx, fmt.Sprintf("%s passed on an empty bag", name))
	} else {
		r.Note(ctx, fmt.Sprintf("%s passed%s", name, chickenFootMessage))
	}

	if round.BaglessPasses >= len(g.Players) {
		g.Note(ctx, "stalemate")
		round.Done = true
		for _, lt := range round.LaidTiles {
			lt.Dead = true
		}
		for _, op := range g.Players {
			op.Dead = true
			op.ChickenFoot = false
		}
	}

	return nil
}

func (g *Game) Note(ctx context.Context, n string) {
	g.History = append(g.History, n)
	log.Print(n)
}

func (g *Game) CurrentRound(ctx context.Context) *Round {
	if len(g.Rounds) == 0 {
		return nil
	}
	r := g.Rounds[len(g.Rounds)-1]
	if r.Done {
		return nil
	}
	return r
}

func (g *Game) DrawTile(ctx context.Context, name string) bool {
	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	if player == nil {
		return false
	}

	if player.JustDrew {
		return false
	}

	if len(g.Bag) > 0 {
		player.Hand = append(player.Hand, g.Bag[0])
		g.Bag = g.Bag[1:]
	}

	player.JustDrew = true

	return true
}

func (g *Game) GetPlayer(ctx context.Context, name string) *Player {
	for _, p := range g.Players {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (g *Game) InBounds(ctx context.Context, c Coord) bool {
	return c.X >= 0 && c.X < g.BoardWidth && c.Y >= 0 && c.Y < g.BoardHeight
}

func (g *Game) LaySpacer(ctx context.Context, name string, spacer *Spacer) error {
	player := g.GetPlayer(ctx, name)
	if player == nil {
		return ErrPlayerNotFound
	}

	if player.ChickenFoot {
		return ErrSpacerNoChickenFoot
	}

	if g.Players[g.Turn].Name != name {
		return ErrNotYourTurn
	}

	round := g.CurrentRound(ctx)
	if round == nil {
		return ErrRoundNotStarted
	}

	round.Spacer = nil

	if err := round.LaySpacer(ctx, g, name, spacer); err != nil {
		return err
	}

	return nil
}

func (g *Game) LayTile(ctx context.Context, name string, tile *LaidTile) error {
	player := g.GetPlayer(ctx, tile.PlayerName)
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

	round := g.CurrentRound(ctx)
	if round == nil {
		return ErrRoundNotStarted
	}

	if g.Players[g.Turn].Name != name {
		return ErrNotYourTurn
	}

	firstTile := len(round.LaidTiles) == 0
	if err := round.LayTile(ctx, g, name, tile, false); err != nil {
		if tile.Indicated != nil && tile.Indicated.PipsA != -1 {
			// try it with the indicated tile
			tile.Indicated = nil
			rt := tile.Reverse()
			if reverseErr := round.LayTile(ctx, g, name, rt, false); reverseErr != nil {
				log.Printf("error with the reverse: %v", reverseErr)
				return err
			}
			tile = rt
			err = nil
		}
		if err != nil {
			return err
		}
	}

	round.Spacer = nil
	if firstTile || tile.Tile.PipsA != tile.Tile.PipsB {
		g.Turn = (g.Turn + 1) % len(g.Players)
	}

	chickenFootMessage := ""
	if player.ChickenFoot {
		player.ChickenFoot = false
		chickenFootMessage = " and is off the foot"
	}
	round.Note(ctx, fmt.Sprintf("%s laid %d:%d%s", name, tile.Tile.PipsA, tile.Tile.PipsB, chickenFootMessage))

	livingPlayers := []*Player{}
	for _, p := range g.Players {
		if !p.Dead {
			livingPlayers = append(livingPlayers, p)
		}
	}

	// playing all tiles immediately wins, even if it causes you to die.
	for _, p := range livingPlayers {
		if len(p.Hand) == 0 {
			round.Done = true
			if len(g.Players) == 1 {
				g.Note(ctx, "you win I guess")
			} else {
				g.Note(ctx, fmt.Sprintf("%s+2 wins the round through efficiency", p.Name))
			}
			p.Score += 2
			for _, lt := range round.LaidTiles {
				if lt.PlayerName == p.Name {
					continue
				}
				lt.Dead = true
			}
			for _, op := range g.Players {
				if op.Name == p.Name {
					continue
				}
				op.Dead = true
				op.ChickenFoot = false
			}
		}
	}
	if !round.Done {
		if len(livingPlayers) == 1 && len(g.Players) > 1 {
			round.Done = true
			g.Note(ctx, fmt.Sprintf("%s+2 wins the round through attrition", livingPlayers[0].Name))
			livingPlayers[0].Score += 2
		} else if len(livingPlayers) == 0 {
			round.Done = true
			if len(g.Players) == 1 {
				g.Note(ctx, "congratulations... you played yourself")
			} else {
				g.Note(ctx, fmt.Sprintf("%s took their ball home", name))
			}
		}
	}

	if round.Done && g.LastRoundLeader(ctx) == 0 {
		g.Done = true
	}

	player.JustDrew = false

	return nil
}

type Player struct {
	Name             string     `json:"name"`
	Score            int        `json:"score"`
	Hand             []*Tile    `json:"hand"`
	Hints            [][]string `json:"hints"`
	SpacerHints      []string   `json:"spacer_hints"`
	ChickenFoot      bool       `json:"chicken_foot"`
	Dead             bool       `json:"dead"`
	JustDrew         bool       `json:"just_drew"`
	ChickenFootCoord Coord      `json:"chicken_foot_coord"`
	ChickenFootURL   string     `json:"chicken_foot_url"`
	Kills            []string   `json:"kills"`
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

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coord) String() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

func (c Coord) Plus(dx, dy int) Coord {
	return Coord{X: c.X + dx, Y: c.Y + dy}
}

func (c Coord) Up() Coord {
	return Coord{X: c.X, Y: c.Y - 1}
}

func (c Coord) Down() Coord {
	return Coord{X: c.X, Y: c.Y + 1}
}

func (c Coord) Left() Coord {
	return Coord{X: c.X - 1, Y: c.Y}
}

func (c Coord) Right() Coord {
	return Coord{X: c.X + 1, Y: c.Y}
}

func (c Coord) OrientationTo(o Coord) string {
	if c.X == o.X {
		if c.Y < o.Y {
			return "down"
		}
		return "up"
	}
	if c.Y == o.Y {
		if c.X < o.X {
			return "right"
		}
		return "left"
	}
	return ""
}

func (c Coord) Adj(o Coord) bool {
	if c.X == o.X {
		return c.Y == o.Y+1 || c.Y == o.Y-1
	}
	if c.Y == o.Y {
		return c.X == o.X+1 || c.X == o.X-1
	}
	return false
}

func (c Coord) Neighbors() []Coord {
	return []Coord{c.Up(), c.Down(), c.Left(), c.Right()}
}

func (c Coord) Orientation(o string) Coord {
	switch o {
	case "up":
		return c.Up()
	case "down":
		return c.Down()
	case "left":
		return c.Left()
	case "right":
		return c.Right()
	}
	return c
}

type Spacer struct {
	A Coord `json:"a"`
	B Coord `json:"b"`
}

type LaidTile struct {
	Tile        *Tile  `json:"tile"`
	Coord       Coord  `json:"coord"`
	Orientation string `json:"orientation"`
	PlayerName  string `json:"player_name"`
	NextPips    int    `json:"next_pips"`
	Dead        bool   `json:"dead"`
	Indicated   *Tile  `json:"indicated"`
}

func (lt *LaidTile) Reverse() *LaidTile {
	rt := &LaidTile{}
	*rt = *lt
	rt.Coord = rt.CoordB()
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

func (lt *LaidTile) CoordA() Coord {
	return lt.Coord
}

func (lt *LaidTile) CoordB() Coord {
	return lt.Coord.Orientation(lt.Orientation)
}

func (lt *LaidTile) String() string {
	return fmt.Sprintf("{%d:%d %s:%s %d}", lt.Tile.PipsA, lt.Tile.PipsB, lt.CoordA(), lt.CoordB(), lt.NextPips)
}

type Round struct {
	Turn          int                    `json:"turn"`
	LaidTiles     []*LaidTile            `json:"laid_tiles"`
	Spacer        *Spacer                `json:"spacer"`
	Done          bool                   `json:"done"`
	History       []string               `json:"history"`
	PlayerLines   map[string][]*LaidTile `json:"player_lines"`
	FreeLines     [][]*LaidTile          `json:"free_lines"`
	BaglessPasses int                    `json:"bagless_passes"`
}

func (r *Round) Note(ctx context.Context, n string) {
	if r == nil {
		return
	}
	r.History = append(r.History, n)
}

func (r *Round) FindHints(ctx context.Context, g *Game, name string, p *Player) {
	squarePips := r.MapTiles(ctx)

	hints := make([]map[string]bool, len(p.Hand))
	for i := range hints {
		hints[i] = map[string]bool{}
	}

	hintAt := func(i int, coord Coord) {
		hints[i][coord.String()] = true
	}

	for i, t := range p.Hand {
		movesOffSquare := func(head *LaidTile, src Coord) {

			for _, orientation := range []string{"up", "down", "left", "right"} {
				lt := &LaidTile{
					Tile:        t,
					Orientation: orientation,
					Coord:       src,
				}
				if r.LayTile(ctx, g, name, lt, true) == nil || r.LayTile(ctx, g, name, lt.Reverse(), true) == nil {
					hintAt(i, lt.CoordA())
					hintAt(i, lt.CoordB())
				}
			}

		}

		movesOffTile := func(head *LaidTile) {
			for _, c := range head.CoordA().Neighbors() {
				movesOffSquare(head, c)
			}
			for _, c := range head.CoordB().Neighbors() {
				movesOffSquare(head, c)
			}
		}

		// first consider direct plays
		for opname, line := range r.PlayerLines {
			op := g.GetPlayer(ctx, opname)
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

		if p.Dead || !p.ChickenFoot {
			for _, line := range r.FreeLines {
				movesOffTile(line[len(line)-1])
			}
		}

		if p.ChickenFoot || r.Spacer == nil {
			// no free line activity allowed
			continue
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

		tryToCoord := func(src Coord) {
			tryA := func(A Coord) {
				for _, orientation := range []string{"up", "down", "left", "right"} {
					lt := &LaidTile{
						Tile:        t,
						Orientation: orientation,
						Coord:       A,
					}
					if r.LayTile(ctx, g, name, lt, true) == nil || r.LayTile(ctx, g, name, lt.Reverse(), true) == nil {
						hintAt(i, lt.CoordA())
						hintAt(i, lt.CoordB())
					}
				}
			}
			tryA(src.Right())
			tryA(src.Left())
			tryA(src.Up())
			tryA(src.Down())
		}
		tryToCoord(r.Spacer.B)
	}
	p.Hints = make([][]string, len(p.Hand))
	for i, hintList := range hints {
		for h := range hintList {
			p.Hints[i] = append(p.Hints[i], h)
		}
	}

	p.SpacerHints = []string{}

	highestLeaderPips := r.PlayerLines[name][0].Tile.PipsA
	for _, line := range r.FreeLines {
		if line[0].Tile.PipsA > highestLeaderPips {
			highestLeaderPips = line[0].Tile.PipsA
		}
	}

	haveAPossibleFreeLine := false
	for _, t := range p.Hand {
		if t.PipsA != t.PipsB {
			continue
		}
		if t.PipsA > highestLeaderPips {
			haveAPossibleFreeLine = true
			break
		}
	}

	if r.Spacer == nil && haveAPossibleFreeLine && !p.ChickenFoot && len(r.PlayerLines[name]) > 1 {
		hintSpacerFrom := func(src Coord) {
			fourWays := []Coord{
				src.Plus(5, 0),
				src.Plus(-5, 0),
				src.Plus(0, 5),
				src.Plus(0, -5),
			}
			for _, dst := range fourWays {
				if g.sixPathFrom(ctx, squarePips, src, dst) {
					p.SpacerHints = append(p.SpacerHints, fmt.Sprintf("%s-%s", src, dst))
				}
			}
		}
		hintSpacerFromTileCoord := func(tc Coord) {
			for _, n := range tc.Neighbors() {
				hintSpacerFrom(n)
			}
		}
		hintSpacerFromTile := func(head *LaidTile) {
			if head.NextPips == head.Tile.PipsA {
				hintSpacerFromTileCoord(head.CoordA())
			}
			if head.NextPips == head.Tile.PipsB {
				hintSpacerFromTileCoord(head.CoordB())
			}
		}
		for _, line := range r.PlayerLines {
			if len(line) == 1 {
				// No spacers off the round leader.
				continue
			}
			hintSpacerFromTile(line[len(line)-1])
		}
		for _, line := range r.FreeLines {
			hintSpacerFromTile(line[len(line)-1])
		}
	}
}

func (r *Round) canPlayOnLine(ctx context.Context, lt *LaidTile, line []*LaidTile) (bool, int, error) {
	last := line[len(line)-1]
	return r.canPlayOnTile(ctx, lt, last)
}

func (r *Round) canPlayOnTile(ctx context.Context, lt, last *LaidTile) (bool, int, error) {
	if lt.Indicated != nil && lt.Indicated.PipsA != -1 {
		if last.Tile.PipsA != lt.Indicated.PipsA || last.Tile.PipsB != lt.Indicated.PipsB {
			return false, 0, ErrMustMatchPips
		}
	}
	return r.canPlayOnTileWithoutIndication(ctx, lt, last)
}

func (r *Round) canPlayOnTileWithoutIndication(ctx context.Context, lt, last *LaidTile) (bool, int, error) {
	var potentialError error
	cerr := func(err error) {
		if err == ErrWrongSide || potentialError == ErrWrongSide {
			potentialError = ErrWrongSide
			return
		}
		if potentialError != nil {
			return
		}
		potentialError = err
	}

	AA := last.CoordA().Adj(lt.CoordA())
	AB := last.CoordA().Adj(lt.CoordB())
	BA := last.CoordB().Adj(lt.CoordA())
	BB := last.CoordB().Adj(lt.CoordB())

	if lt.Tile.PipsA == last.NextPips {
		if last.Tile.PipsA == lt.Tile.PipsA {
			if AA {
				return true, lt.Tile.PipsB, nil
			}
			if AB {
				cerr(ErrWrongSide)
			}
			cerr(ErrNotAdjacent)
		}
		if last.Tile.PipsB == lt.Tile.PipsA {
			if BA {
				return true, lt.Tile.PipsB, nil
			}
			if BB {
				cerr(ErrWrongSide)
			}
			cerr(ErrNotAdjacent)
		}
	}
	if lt.Tile.PipsB == last.NextPips {
		if last.Tile.PipsA == lt.Tile.PipsB {
			if AB {
				return true, lt.Tile.PipsA, nil
			}
			if AA {
				cerr(ErrWrongSide)
			}
			cerr(ErrNotAdjacent)
		}
		if last.Tile.PipsB == lt.Tile.PipsB {
			if BB {
				return true, lt.Tile.PipsA, nil
			}
			if BA {
				cerr(ErrWrongSide)
			}
			cerr(ErrNotAdjacent)
		}
	}
	cerr(ErrMustMatchPips)
	return false, 0, potentialError
}

func (g *Game) sixPathFrom(ctx context.Context, squarePips map[Coord]SquarePips, src, dst Coord) bool {
	if !g.InBounds(ctx, src) {
		return false
	}
	if !g.InBounds(ctx, dst) {
		return false
	}

	dx := dst.X - src.X
	dy := dst.Y - src.Y

	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}

	if dx == 0 && dy != 5 {
		return false
	}
	if dy == 0 && dx != 5 {
		return false
	}

	orientation := src.OrientationTo(dst)
	cur := src

	check := func(c Coord) bool {
		if _, ok := squarePips[c]; ok {
			return false
		}
		return true
	}

	for i := 0; i < 5; i++ {
		if !check(cur) {
			return false
		}
		cur = cur.Orientation(orientation)
	}
	return check(cur)
}

func (r *Round) LaySpacer(ctx context.Context, g *Game, name string, spacer *Spacer) error {
	if r.Done {
		return ErrRoundAlreadyDone
	}

	if (*spacer) == (Spacer{}) {
		r.Spacer = nil
		return nil
	}

	if spacer.A.X != spacer.B.X && spacer.A.Y != spacer.B.Y {
		return ErrSpacerNotStraight
	}
	switch 5 {
	case spacer.A.X - spacer.B.X:
	case spacer.B.X - spacer.A.X:
	case spacer.A.Y - spacer.B.Y:
	case spacer.B.Y - spacer.A.Y:
	default:
		return ErrWrongLengthSpacer
	}

	if !g.sixPathFrom(ctx, r.MapTiles(ctx), spacer.A, spacer.B) {
		return ErrTileOccluded
	}

	// verify that x1,y1 is adjacent to a line head.
	checkLineHead := func(lt *LaidTile) bool {
		if lt.NextPips == lt.Tile.PipsA {
			return spacer.A.Adj(lt.CoordA())
		}
		if lt.NextPips == lt.Tile.PipsB {
			return spacer.A.Adj(lt.CoordB())
		}
		return false
	}
	isOnLineHead := false
	for _, line := range r.PlayerLines {
		if checkLineHead(line[len(line)-1]) {
			isOnLineHead = true
			break
		}
	}
	for _, line := range r.FreeLines {
		if checkLineHead(line[len(line)-1]) {
			isOnLineHead = true
			break
		}
	}
	if !isOnLineHead {
		return ErrSpacerNotStartedOnLine
	}

	r.Spacer = spacer
	return nil
}

func (r *Round) LayTile(ctx context.Context, g *Game, name string, lt *LaidTile, dryRun bool) error {
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

	if !g.InBounds(ctx, lt.CoordA()) {
		return ErrTileOutOfBounds
	}
	if !g.InBounds(ctx, lt.CoordB()) {
		return ErrTileOutOfBounds
	}

	squarePips := r.MapTiles(ctx)
	if _, ok := squarePips[lt.CoordA()]; ok {
		return ErrTileOccluded
	}
	if _, ok := squarePips[lt.CoordB()]; ok {
		return ErrTileOccluded
	}

	if r.BlockingFeet(ctx, g, squarePips, lt, name) {
		return ErrNoBlockingFeet
	}

	playerFoot := ""
	// if this is on someone's foot, it's definitely made part of their line.
	for _, p := range g.Players {
		if p.ChickenFootCoord == lt.CoordA() {
			playerFoot = p.Name
		}
		if p.ChickenFootCoord == lt.CoordB() {
			playerFoot = p.Name
		}
	}

	var potentialError error
	cerr := func(err error) {
		precedence := []error{ErrWrongSide, ErrMustMatchPips, ErrNotAdjacent}
		for _, e := range precedence {
			if potentialError == e || err == e {
				potentialError = err
				return
			}
		}
		potentialError = err
	}

	playedALine := false

	player := g.GetPlayer(ctx, name)

	if r.Spacer == nil && !player.Dead && (playerFoot == "" || playerFoot == player.Name) {
		mainLine := r.PlayerLines[player.Name]

		onFoot := false
		if player.ChickenFoot && len(mainLine) == 1 {
			if player.ChickenFootCoord == lt.CoordA() {
				onFoot = true
			}
			if player.ChickenFootCoord == lt.CoordB() {
				onFoot = true
			}
		}
		if len(mainLine) == 1 && !onFoot && player.ChickenFoot {
			return ErrMustPlayOnFoot
		}
		if len(mainLine) > 1 || onFoot || !player.ChickenFoot {
			if ok, nextPips, err := r.canPlayOnLine(ctx, lt, mainLine); ok {
				playedALine = true
				if !dryRun {
					r.PlayerLines[player.Name] = append(mainLine, lt)
					lt.NextPips = nextPips
				}
			} else {
				cerr(err)
			}
		}
	}

	canPlayOtherLines := r.Spacer == nil && (player.Dead || !player.ChickenFoot)

	if canPlayOtherLines {
		for oname, line := range r.PlayerLines {
			if playerFoot != "" && playerFoot != oname {
				continue
			}
			if playedALine {
				continue
			}
			op := g.GetPlayer(ctx, oname)
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
				if lt.CoordA() == op.ChickenFootCoord && lt.Tile.PipsA == line[0].NextPips {
					isOnMyFoot = true
				}
				if lt.CoordB() == op.ChickenFootCoord && lt.Tile.PipsB == line[0].NextPips {
					isOnMyFoot = true
				}
				if !isOnMyFoot {
					continue
				}
			}
			if ok, nextPips, err := r.canPlayOnLine(ctx, lt, line); ok {
				playedALine = true
				if !dryRun {
					lt.PlayerName = oname
					r.PlayerLines[oname] = append(line, lt)
					lt.NextPips = nextPips
					// Update the chicken-foot.
					if lt.NextPips == lt.Tile.PipsA {
						op.ChickenFootCoord = lt.CoordB()
					}
					if lt.NextPips == lt.Tile.PipsB {
						op.ChickenFootCoord = lt.CoordA()
					}
				}
			} else {
				cerr(err)
			}
		}
		for i, line := range r.FreeLines {
			if playedALine {
				continue
			}
			if ok, nextPips, err := r.canPlayOnLine(ctx, lt, line); ok {
				playedALine = true
				if !dryRun {
					lt.PlayerName = ""
					r.FreeLines[i] = append(line, lt)
					lt.NextPips = nextPips
				}
			} else {
				cerr(err)
			}
		}
	}

	canStartFreeLine := r.Spacer != nil && (player.Dead || !player.ChickenFoot)
	playedAtLeastOne := len(r.PlayerLines[player.Name]) > 1
	isDouble := lt.Tile.PipsA == lt.Tile.PipsB
	if canStartFreeLine && isDouble && playedAtLeastOne {
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
		if isHigher {

			inSpacer := func(src Coord) bool {
				if src.X >= r.Spacer.A.X && src.X <= r.Spacer.B.X && src.Y >= r.Spacer.A.Y && src.Y <= r.Spacer.B.Y {
					return true
				}
				return false
			}
			canBeFree := r.Spacer.B.Adj(lt.CoordA()) || r.Spacer.B.Adj(lt.CoordB())
			if inSpacer(lt.CoordA()) || inSpacer(lt.CoordB()) {
				canBeFree = false
			}

			if canBeFree {
				playedALine = true
				if !dryRun {
					lt.PlayerName = ""
					r.FreeLines = append(r.FreeLines, []*LaidTile{lt})
					lt.NextPips = lt.Tile.PipsA
					g.Note(ctx, fmt.Sprintf("%s started a free line", name))
				}
			}
		}
	}

	if !playedALine {
		if potentialError != nil {
			return potentialError
		}
		return ErrNoLine
	}

	// Add the new tile to the occlusion grid.
	squarePips[lt.CoordA()] = SquarePips{LaidTile: lt, Pips: lt.Tile.PipsA}
	squarePips[lt.CoordB()] = SquarePips{LaidTile: lt, Pips: lt.Tile.PipsB}
	isOpenFrom := func(src Coord) bool {
		for _, na := range src.Neighbors() {
			if !g.InBounds(ctx, na) {
				continue
			}
			if _, ok := squarePips[na]; ok {
				continue
			}
			for _, nb := range na.Neighbors() {
				if !g.InBounds(ctx, nb) {
					continue
				}
				if _, ok := squarePips[nb]; ok {
					continue
				}
				return true
			}
		}
		return false
	}
	isCutOff := func(line []*LaidTile) bool {
		last := line[len(line)-1]
		if last.NextPips == last.Tile.PipsA && isOpenFrom(last.CoordA()) {
			return false
		}
		if last.NextPips == last.Tile.PipsB && isOpenFrom(last.CoordB()) {
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
				g.Note(ctx, fmt.Sprintf("what kind of reprobate cuts off a free line? (%s+0)", name))
				tilesInFreeLine := map[string]bool{}
				for _, lt := range line {
					tilesInFreeLine[lt.Tile.String()] = true
				}
				for _, lt := range r.LaidTiles {
					if tilesInFreeLine[lt.Tile.String()] {
						lt.Dead = true
					}
				}
				lt.Dead = true
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
				g.Note(ctx, fmt.Sprintf("%s+1 ducked out of that situation", player.Name))
			} else {
				g.Note(ctx, fmt.Sprintf("%s+1 cut-off %s's-1 line", player.Name, p.Name))
			}
			player.Score += 1
			player.Kills = append(player.Kills, p.Name)
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
	if !dryRun && canPlayOtherLines {
		consumed := []string{}

		canConsume := func(head *LaidTile) bool {
			if lt.NextPips != head.NextPips {
				return false
			}
			checkLTCoord := func(src Coord) bool {
				if head.Tile.PipsA == head.NextPips {
					if src.Adj(head.CoordA()) {
						return true
					}
				}
				if head.Tile.PipsB == head.NextPips {
					if src.Adj(head.CoordB()) {
						return true
					}
				}
				return false
			}
			if lt.Tile.PipsA == lt.NextPips {
				if checkLTCoord(lt.CoordA()) {
					return true
				}
			}
			if lt.Tile.PipsB == lt.NextPips {
				if checkLTCoord(lt.CoordB()) {
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
				consumed = append(consumed, fmt.Sprintf("%s-1", lt.PlayerName))
			}
			freeLineNote := ""
			if freeLines > 0 {
				freeLineNote = fmt.Sprintf("%d free line(s)", freeLines)
			}
			g.Note(ctx, fmt.Sprintf("%s's+%d IL OUROBOROS consumes %s", name, len(consumed), strings.Join(append(consumed, freeLineNote), ", ")))
			for _, n := range consumed {
				op := g.GetPlayer(ctx, n)
				op.Dead = true
				op.ChickenFoot = false
				for _, lt := range r.PlayerLines[op.Name][1:] {
					lt.Dead = true
					deadTileKeys[lt.Tile.String()] = true
				}
				op.Score -= 1
				player.Score += 1
				player.Kills = append(player.Kills, op.Name)
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

func (r *Round) BlockingFeet(ctx context.Context, g *Game, squarePips map[Coord]SquarePips, ot *LaidTile, name string) bool {
	lt := &LaidTile{}
	*lt = *ot
	lt.PlayerName = name

	allFrom := func(src Coord, orientations []string) []*LaidTile {
		lts := []*LaidTile{}
		for _, o := range orientations {
			lts = append(lts, &LaidTile{
				Coord:       src,
				Orientation: o,
			})
		}
		return lts
	}

	blocks := map[Coord]bool{}
	for coord := range squarePips {
		blocks[coord] = true
	}

	playersToSatisfy := []*Player{}
	playerChickenFeetCoords := map[string]Coord{}
	// Does this set of blocks prevent any player from starting their main line?
	for pn, pline := range r.PlayerLines {
		if lt.PlayerName == pn {
			continue
		}
		if len(pline) > 1 {
			continue
		}
		p := g.GetPlayer(ctx, pn)
		if p.ChickenFoot {
			playerChickenFeetCoords[p.Name] = p.ChickenFootCoord
			if p.ChickenFootCoord == lt.CoordA() || p.ChickenFootCoord == lt.CoordB() {
				lt.PlayerName = p.Name
				return false
			}
		}
		playersToSatisfy = append(playersToSatisfy, p)
	}

	// depth := 0
	// l := func(fmt string, items ...any) {
	// 	prefix := strings.Repeat(" ", depth)
	// 	log.Printf(prefix+fmt, items...)
	// }

	// returns true if the player can be satisfied.
	var recursiveEnsurePlayersOK func(playersLeft []*Player, prevBlocks map[Coord]bool, lt *LaidTile) bool
	recursiveEnsurePlayersOK = func(playersLeft []*Player, prevBlocks map[Coord]bool, lt *LaidTile) bool {
		// l("dropping %s:%s for %s", lt.CoordA(), lt.CoordB(), lt.PlayerName)

		// depth += 1
		// defer func() { depth -= 1 }()

		if len(playersLeft) == 0 {
			// l("looks good")
			return true
		}

		p := playersLeft[0]
		// l("considering player %s", p.Name)
		// Definitely not on someone else's reserved foot.
		for op, coord := range playerChickenFeetCoords {
			if op == lt.PlayerName {
				continue
			}
			if coord == lt.CoordA() || coord == lt.CoordB() {
				// l("on %s's foot", op)
				return false
			}
		}

		blocks := map[Coord]bool{}
		for coord := range prevBlocks {
			blocks[coord] = true
		}
		blocks[lt.CoordA()] = true
		blocks[lt.CoordB()] = true

		isOpen := func(src Coord) bool {
			if playerChickenFeetCoords[p.Name] == src {
				return true
			}
			return !blocks[src]
		}
		checkCoord := func(src Coord, orientations []string) (bool, bool) {
			// l("checking coord A %d,%d", x, y)
			if !isOpen(src) {
				// l("%s A-blocked at %d,%d", p.Name, x, y)
				return false, false
			}
			possibilities := allFrom(src, orientations)
			canFitMyself := false
			for _, pos := range possibilities {
				pos.PlayerName = p.Name
				// l("checking coord B %s", pos.CoordB())
				if !isOpen(pos.CoordB()) {
					// l("%s B-blocked at %s:%s", p.Name, pos.CoordA(), pos.CoordB())
					continue
				}
				canFitMyself = true
				// l("trying %s:%s with %v", pos.CoordA(), pos.CoordB(), playersLeft[1:])
				if recursiveEnsurePlayersOK(playersLeft[1:], blocks, pos) {
					return true, canFitMyself
				}
			}
			return false, canFitMyself
		}

		if p.ChickenFoot {
			success, _ := checkCoord(p.ChickenFootCoord, []string{"up", "down", "left", "right"})
			// We ignore the canFitMyself return value because this is the only place they can try.
			return success
		}

		// iterate around the round leader
		rl := r.PlayerLines[p.Name][0]
		var success, canFitMyself bool
		if success, canFitMyself = checkCoord(rl.CoordA().Left(), []string{"left", "up", "down"}); success {
			return true
		}
		// we stop early, no point in checking further since this tile could just swap with
		// other unplayed lines. This prevents a lot of branching.
		if canFitMyself {
			return false
		}
		if success, canFitMyself = checkCoord(rl.CoordA().Down(), []string{"down", "left", "right"}); success {
			return true
		}
		if canFitMyself {
			return false
		}
		if success, canFitMyself = checkCoord(rl.CoordA().Up(), []string{"up", "left", "right"}); success {
			return true
		}
		if canFitMyself {
			return false
		}
		if success, canFitMyself = checkCoord(rl.CoordB().Right(), []string{"right", "up", "down"}); success {
			return true
		}
		if canFitMyself {
			return false
		}
		if success, canFitMyself = checkCoord(rl.CoordB().Down(), []string{"down", "left", "right"}); success {
			return true
		}
		if canFitMyself {
			return false
		}
		if success, _ = checkCoord(rl.CoordB().Up(), []string{"up", "left", "right"}); success {
			return true
		}
		return false
	}

	return !recursiveEnsurePlayersOK(playersToSatisfy, blocks, lt)
}

type SquarePips struct {
	LaidTile *LaidTile
	Pips     int
}

func (r *Round) MapTiles(ctx context.Context) map[Coord]SquarePips {
	squarePips := map[Coord]SquarePips{}
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
