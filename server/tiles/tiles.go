package tiles

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/golang/protobuf/proto"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

const (
	handSize = 10
)

var Debug = false

func debug(format string, items ...interface{}) {
	if !Debug {
		return
	}
	fmt.Printf(format, items...)
}

func SetupBoard(ctx context.Context, b *tpb.Board, lastLeaderPips int32) (*tpb.Board, error) {
	if lastLeaderPips == 0 {
		// First round, highest available double leads.
		lastLeaderPips = 10000
	}
	if lastLeaderPips < 2 {
		return nil, errors.New("last leader must be at least 2")
	}
	if len(b.GetPlayers()) > 6 {
		return nil, errors.New("only 6 players at most")
	}
	if len(b.GetPlayers()) < 2 {
		return nil, errors.New("at least 2 players required")
	}

	// In tests, the bag is provided so we don't shuffle a new one.
	if b.Bag == nil {
		b.Bag = newShuffledBag(ctx)
	}

	if handSize*len(b.GetPlayers()) > len(b.GetBag()) {
		return nil, fmt.Errorf("%d players need %d tiles, only had %d", len(b.GetPlayers()), handSize*len(b.GetPlayers()), len(b.GetBag()))
	}

	var firstPlayer *tpb.Player
	var firstTile *tpb.Tile

	// Each player gets their initial draw.
	for _, p := range b.GetPlayers() {
		p.Hand = b.Bag[:handSize]
		b.Bag = b.Bag[handSize:]
	}
	for firstTile == nil {
		for _, p := range b.GetPlayers() {
			// Go through the tiles drawn and check if any are a candidate for the round leader.
			for _, t := range p.Hand {
				// Only doubles.
				if t.GetA() != t.GetB() {
					continue
				}
				// Has to be lower than the last leader.
				if t.GetA() >= lastLeaderPips {
					continue
				}
				// Among the remaining, we choose the highest.
				if firstTile != nil && t.GetA() < firstTile.GetA() {
					continue
				}
				firstTile = t
				firstPlayer = p
			}
		}
		if firstTile == nil {
			// No leader found, everyone gets another tile and we dry again.
			for _, p := range b.GetPlayers() {
				p.Hand = append(p.Hand, b.Bag[0])
				b.Bag = b.Bag[1:]
			}
		}
	}

	// Remove the new round leader from the player's hand.
	fphand, ok := removeTileFromHand(ctx, firstTile, firstPlayer.GetHand())
	if !ok {
		return nil, fmt.Errorf("could not remove %v from %v", firstTile, firstPlayer.GetHand())
	}
	firstPlayer.Hand = fphand
	var err error
	b.NextPlayerId, err = nextPlayer(ctx, b, firstPlayer.GetPlayerId())
	if err != nil {
		return nil, fmt.Errorf("could not determine next player: %v", err)
	}

	// Create a new line for each player, beginning with the round leader in the middle of the
	// board as the first placement.

	leaderPlacement := &tpb.Placement{
		Tile: firstTile,
		A: &tpb.Coord{
			X: b.GetWidth()/2 - 1,
			Y: b.GetHeight() / 2,
		},
		B: &tpb.Coord{
			X: b.GetWidth() / 2,
			Y: b.GetHeight() / 2,
		},
		Type: tpb.Placement_PLAYER_LEADER,
	}

	for _, p := range b.GetPlayers() {
		pl := &tpb.Line{
			Placements: []*tpb.Placement{leaderPlacement},
			PlayerId:   p.GetPlayerId(),
		}
		b.PlayerLines = append(b.PlayerLines, pl)
	}

	return b, nil
}

func nextPlayer(ctx context.Context, b *tpb.Board, playerID string) (string, error) {
	firstPlayer := b.GetPlayers()[0].GetPlayerId()
	for i, p := range b.GetPlayers() {
		if playerID != p.GetPlayerId() {
			continue
		}
		if i == len(b.GetPlayers())-1 {
			return firstPlayer, nil
		}
		return b.GetPlayers()[i+1].GetPlayerId(), nil
	}
	return "", fmt.Errorf("player %q not found", playerID)
}

func removeTileFromHand(ctx context.Context, t *tpb.Tile, hand []*tpb.Tile) ([]*tpb.Tile, bool) {
	for i, ht := range hand {
		if t != ht {
			continue
		}
		nh := append([]*tpb.Tile{}, hand[:i]...)
		nh = append(nh, hand[i+1:]...)
		return nh, true
	}
	return hand, false
}

func newShuffledBag(ctx context.Context) []*tpb.Tile {
	b := []*tpb.Tile{}
	for i := 0; i <= 12; i++ {
		for j := i; j <= 12; j++ {
			b = append(b, &tpb.Tile{
				A: int32(i),
				B: int32(j),
			})
		}
	}
	rand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return b
}

func LayTile(ctx context.Context, b *tpb.Board, move *tpb.Placement) (*tpb.Board, error) {
	players := map[string]*tpb.Player{}
	for _, p := range b.GetPlayers() {
		players[p.GetPlayerId()] = p
	}
	player, ok := players[b.GetNextPlayerId()]
	if !ok {
		return nil, fmt.Errorf("next player %s not in list", b.GetNextPlayerId())
	}

	validMoves, err := LegalMoves(ctx, b, player)
	if err != nil {
		return nil, fmt.Errorf("could not get legal moves before playing: %v", err)
	}
	legal := false
	debug("checking %s\n", move)
	for _, vm := range validMoves {
		debug("is it %s?\n", vm)
		if !proto.Equal(vm.GetTile(), move.GetTile()) {
			debug("no\n")
			continue
		}
		if !proto.Equal(vm.GetA(), move.GetA()) {
			debug("no\n")
			continue
		}
		if !proto.Equal(vm.GetB(), move.GetB()) {
			debug("no\n")
			continue
		}
		debug("yes\n")
		move.Type = vm.GetType()
		legal = true
	}
	if !legal {
		return nil, fmt.Errorf("%v was an illegal move", move)
	}

	// If the player passed, apply chickenfooting and move to the next turn.
	if move.Type == tpb.Placement_PASS {
		player.ChickenFooted = true
		b.NextPlayerId, err = nextPlayer(ctx, b, b.GetNextPlayerId())
		if err != nil {
			return nil, fmt.Errorf("could not get next player: %v", err)
		}
		debug("passed, turn goes to %s\n", b.GetNextPlayerId())
		b.Done, err = gameIsOver(ctx, b)
		if err != nil {
			return nil, fmt.Errorf("could not check if game was done after pass: %v", err)
		}
		return b, nil
	}

	// If the player draws, add the next tile from the bag to their hand
	// and move to the next turn.
	if move.Type == tpb.Placement_DRAW {
		if len(b.GetBag()) == 0 {
			return nil, errors.New("tried to draw when there were no tiles in the bag")
		}
		drawnTile := b.GetBag()[0]
		b.Bag = b.GetBag()[1:]
		player.Hand = append(player.Hand, drawnTile)
		b.NextPlayerId, err = nextPlayer(ctx, b, b.GetNextPlayerId())
		if err != nil {
			return nil, fmt.Errorf("could not get next player: %v", err)
		}
		debug("drew, turn goes to %s\n", b.GetNextPlayerId())
		b.Done, err = gameIsOver(ctx, b)
		if err != nil {
			return nil, fmt.Errorf("could not check if game was done after draw: %v", err)
		}
		return b, nil
	}

	debug("must be a continuation\n")

	// Check all lines to see if this tile can be placed on them.
	availableLines, err := AvailableLines(ctx, b, player, players)
	if err != nil {
		return nil, fmt.Errorf("could not get available lines: %v", err)
	}
	playableLines := []*tpb.Line{}
	for _, l := range availableLines {
		starts, pips, err := NextStartsAndPips(ctx, l)
		if err != nil {
			return nil, fmt.Errorf("could not get starts and pips for line: %v", err)
		}
		if pips == move.GetTile().GetA() && isInCoordList(ctx, move.GetA(), starts) {
			playableLines = append(playableLines, l)
			continue
		}
		if pips == move.GetTile().GetB() && isInCoordList(ctx, move.GetB(), starts) {
			playableLines = append(playableLines, l)
		}
	}

	if len(playableLines) == 0 {
		return nil, errors.New("illegal move")
	}

	debug("%d playable lines\n", len(playableLines))

	for _, l := range playableLines {
		// If a move can be played on more than one line, they all die (il ouroboros).
		if len(playableLines) > 1 {
			l.Murderer = player.GetPlayerId()
		}
		l.Placements = append(l.Placements, move)
	}

	// Check all lines to see if there is room for another play sometime
	// in the future. If not, the line is dead.
	m := maskForBoard(ctx, b)
	for _, l := range b.GetPlayerLines() {
		if l.GetMurderer() != "" {
			continue
		}
		room, err := roomToPlay(ctx, m, l)
		if err != nil {
			return nil, fmt.Errorf("could not check for room: %v", err)
		}
		if !room {
			l.Murderer = player.GetPlayerId()
		}
	}

	// Remove the laid tile from the player's hand.
	nh := []*tpb.Tile{}
	for _, t := range player.GetHand() {
		if proto.Equal(t, move) {
			continue
		}
		nh = append(nh, t)
	}
	player.Hand = nh

	// Once you play a tile you're no longer chickenfooted.
	player.ChickenFooted = false

	over, err := gameIsOver(ctx, b)
	if err != nil {
		return nil, fmt.Errorf("could not check if game is over: %v", err)
	}
	if over {
		b.Done = true
		debug("game over\n")
		return b, nil
	}

	b.NextPlayerId, err = nextPlayer(ctx, b, b.GetNextPlayerId())
	if err != nil {
		return nil, fmt.Errorf("could not get next player: %v", err)
	}

	debug("laid, turn goes to %s\n", b.GetNextPlayerId())
	return b, nil
}

func GetNextPlayer(ctx context.Context, b *tpb.Board) (*tpb.Player, error) {
	for _, p := range b.GetPlayers() {
		if p.GetPlayerId() == b.GetNextPlayerId() {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no player with id %q", b.GetNextPlayerId())
}

func gameIsOver(ctx context.Context, b *tpb.Board) (bool, error) {

	for _, player := range b.GetPlayers() {
		if len(player.GetHand()) == 0 {
			debug("game is over because a %s ran out of tiles\n", player.GetPlayerId())
			return true, nil
		}
	}

	linesAlive := 0
	for _, pl := range b.GetPlayerLines() {
		if pl.GetMurderer() == "" {
			linesAlive++
		}
	}
	if linesAlive < 2 {
		debug("game is over because there aren't at least two player lines left\n")
		return true, nil
	}

	playersWithMoves := 0
	for _, player := range b.GetPlayers() {
		moves, err := LegalMoves(ctx, b, player)
		if err != nil {
			return false, fmt.Errorf("could not find legal moves for %s: %v", player.GetPlayerId(), err)
		}
		if len(moves) == 1 && moves[0].GetType() == tpb.Placement_PASS {
			continue
		}
		fmt.Printf("internal: %s has %d moves\n", player.GetPlayerId(), len(moves))
		playersWithMoves++
	}
	if playersWithMoves == 0 {
		debug("game is over because no player has any moves\n")
		return true, nil
	}

	return false, nil
}

func roomToPlay(ctx context.Context, m mask, line *tpb.Line) (bool, error) {
	starts, _, err := NextStartsAndPips(ctx, line)
	if err != nil {
		return false, fmt.Errorf("could not get starts: %v", err)
	}
	for _, start := range starts {
		if !m.getc(start) {
			continue
		}
		ends := getAdjacent(ctx, start)
		for _, end := range ends {
			if m.getc(end) {
				return true, nil
			}
		}
	}
	return false, nil
}

func isInCoordList(ctx context.Context, c *tpb.Coord, cs []*tpb.Coord) bool {
	for _, o := range cs {
		if proto.Equal(c, o) {
			return true
		}
	}
	return false
}

func AvailableLines(ctx context.Context, b *tpb.Board, player *tpb.Player, players map[string]*tpb.Player) ([]*tpb.Line, error) {
	// Check the lines that can be played on. The player's own line is always included. If
	// the player is not chickenfooted, then other chickenfooted lines are also included.
	availableLines := []*tpb.Line{}
	for _, l := range b.GetPlayerLines() {
		if l.GetPlayerId() == player.GetPlayerId() {
			availableLines = append(availableLines, l)
			continue
		}

		// If the player is chickenfooted, the only legal move is on their own line.
		if player.GetChickenFooted() {
			continue
		}

		lp, ok := players[l.GetPlayerId()]
		if !ok {
			return nil, fmt.Errorf("bad board state, line has unknown player %q", l.GetPlayerId())
		}
		if lp.GetChickenFooted() {
			availableLines = append(availableLines, l)
		}
	}

	// If the player is not chickenfooted, they can also play on any free lines.
	if !player.GetChickenFooted() {
		availableLines = append(availableLines, b.GetFreeLines()...)
	}
	return availableLines, nil
}

func NextStartsAndPips(ctx context.Context, line *tpb.Line) ([]*tpb.Coord, int32, error) {
	placements := line.GetPlacements()
	if len(placements) == 0 {
		return nil, 0, errors.New("line had no placements")
	}

	lastPlacement := placements[len(placements)-1]
	if lastPlacement.GetTile().GetA() == lastPlacement.GetTile().GetB() {
		var starts []*tpb.Coord
		starts = append(starts, getAdjacent(ctx, lastPlacement.GetA())...)
		starts = append(starts, getAdjacent(ctx, lastPlacement.GetB())...)
		pips := lastPlacement.GetTile().GetA()
		return starts, pips, nil
	}

	pips := placements[0].GetTile().GetA() // It's necessarily a double.
	var curSide *tpb.Coord
	for _, placement := range placements[1:] {
		if placement.GetTile().GetA() == pips {
			pips = placement.GetTile().GetB()
			curSide = placement.GetB()
			continue
		}
		if placement.GetTile().GetB() == pips {
			pips = placement.GetTile().GetA()
			curSide = placement.GetA()
			continue
		}
		return nil, 0, fmt.Errorf("problem with this line at %v", placement)
	}
	starts := getAdjacent(ctx, curSide)
	return starts, pips, nil
}

func LegalMoves(ctx context.Context, b *tpb.Board, p *tpb.Player) ([]*tpb.Placement, error) {
	if b.Done {
		return nil, nil
	}

	// Quick lookup to map IDs to players later, mostly for seeing if a player is chickenfooted.
	players := map[string]*tpb.Player{}
	for _, p := range b.GetPlayers() {
		players[p.GetPlayerId()] = p
	}

	moves := []*tpb.Placement{}

	// Check the lines that can be played on. The player's own line is always included. If
	// the player is not chickenfooted, then other chickenfooted lines are also included.
	availableLines, err := AvailableLines(ctx, b, p, players)
	if err != nil {
		return nil, fmt.Errorf("could not get available lines: %v", err)
	}

	debug("%s wants to play\n", p.GetPlayerId())
	debug("lines available: %d\n", len(availableLines))

	// Build a mask of the board so we can easily tell where a tile can be placed.
	mask := maskForBoard(ctx, b)

	// For each available line, look at the last placement. We compare it to each tile in the
	// player's hand, and then check the ways that tile could be placed that don't hit the
	// mask.
	for i, l := range availableLines {
		// Can't play on a dead line.
		if l.GetMurderer() != "" {
			continue
		}
		starts, pips, err := NextStartsAndPips(ctx, l)
		if err != nil {
			return nil, fmt.Errorf("could not get next starts and pips: %v", err)
		}

		var placementType tpb.Placement_Type
		switch l.GetPlacements()[0].GetType() {
		case tpb.Placement_PLAYER_LEADER:
			placementType = tpb.Placement_PLAYER_CONTINUATION
		case tpb.Placement_FREE_LEADER:
			placementType = tpb.Placement_FREE_CONTINUATION
		default:
			return nil, fmt.Errorf("line start wasn't a leader, instead %s", l.GetPlacements()[0].GetType().String())
		}

		debug("Line %d can start with %d pips at any of %q\n", i, pips, starts)

		debug("Player %s has hand %q\n", p.GetPlayerId(), p.GetHand())

		// We consider any tile that can place `pips` in one of the `starts`.
		for _, t := range p.GetHand() {
			// Unfortunate duplication.
			debug("Looking at tile %q\n", t)
			if t.GetA() == pips {
				// debug("%q matches on the A side\n", t)
				for _, a := range starts {
					// debug("Checking when A=%q\n", a)
					bs := getAdjacent(ctx, a)
					for _, b := range bs {
						nextPlacement := &tpb.Placement{
							Tile: t,
							A:    a,
							B:    b,
							Type: placementType,
						}
						if !mask.getp(nextPlacement) {
							// debug("%q doesn't hit the mask\n", nextPlacement)
							moves = append(moves, nextPlacement)
						}
						// debug("%q hits the mask\n", nextPlacement)
					}
				}
			}
			if t.GetB() == pips {
				// debug("%q matches on the B side\n", t)
				for _, b := range starts {
					// debug("Checking when B=%q\n", b)
					as := getAdjacent(ctx, b)
					for _, a := range as {
						nextPlacement := &tpb.Placement{
							Tile: t,
							A:    a,
							B:    b,
							Type: placementType,
						}
						if !mask.getp(nextPlacement) {
							// debug("%q doesn't hit the mask\n", nextPlacement)
							moves = append(moves, nextPlacement)
						}
						// debug("%q hits the mask\n", nextPlacement)
					}
				}
			}
		}
	}

	moves = append(moves, &tpb.Placement{
		Type: tpb.Placement_PASS,
	})
	if len(b.GetBag()) != 0 {
		moves = append(moves, &tpb.Placement{
			Type: tpb.Placement_DRAW,
		})
	}

	// TODO: opportunities to start new free lines.

	// TODO: disallow moves that "crowd" the round leader.

	return moves, nil
}

func getAdjacent(ctx context.Context, c *tpb.Coord) []*tpb.Coord {
	return []*tpb.Coord{{
		X: c.X,
		Y: c.Y + 1,
	}, {
		X: c.X,
		Y: c.Y - 1,
	}, {
		X: c.X + 1,
		Y: c.Y,
	}, {
		X: c.X - 1,
		Y: c.Y,
	}}
}
