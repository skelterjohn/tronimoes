package tiles

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

const (
	handSize = 10
)

func SetupBoard(ctx context.Context, b *tpb.Board, lastLeaderPips int32) (*tpb.Board, error) {
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

func LegalMoves(ctx context.Context, b *tpb.Board, playerID string) ([]*tpb.Placement, error) {
	// Quick lookup to map IDs to players later, mostly for seeing if a player is chickenfooted.
	players := map[string]*tpb.Player{}
	for _, p := range b.GetPlayers() {
		players[p.GetPlayerId()] = p
	}

	p, ok := players[playerID]
	if !ok {
		return nil, fmt.Errorf("no player %q", playerID)
	}

	moves := []*tpb.Placement{}

	// Check the lines that can be played on. The player's own line is always included. If
	// the player is not chickenfooted, then other chickenfooted lines are also included.
	availableLines := []*tpb.Line{}
	for _, l := range b.GetPlayerLines() {
		if l.GetPlayerId() == p.GetPlayerId() {
			availableLines = append(availableLines, l)
			continue
		}

		// If the player is chickenfooted, the only legal move is on their own line.
		if p.GetChickenFooted() {
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
	if !p.GetChickenFooted() {
		availableLines = append(availableLines, b.GetFreeLines()...)
	}

	// Build a mask of the board so we can easily tell where a tile can be placed.
	mask := maskForBoard(ctx, b)

	// For each available line, look at the last placement. We compare it to each tile in the
	// player's hand, and then check the ways that tile could be placed that don't hit the
	// mask.
	for _, l := range availableLines {
		placements := l.GetPlacements()
		lp := placements[len(placements)-1]

		starts := []*tpb.Coord{}
		var pips int32 = 0

		if lp.GetTile().GetA() == lp.GetTile().GetB() {
			// If the tile is a double, we can start from either side with the next tile.
			starts = append(starts, getAdjacent(ctx, lp.GetA())...)
			starts = append(starts, getAdjacent(ctx, lp.GetB())...)
			pips = lp.GetTile().GetA()
		} else {
			// Otherwise, start from the beginning and follow the line until the end
			// to be sure we have the right pips and start point.
			pips = placements[0].GetTile().GetA() // It's necessarily a double.
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
				return nil, fmt.Errorf("problem with this line at %v", placement)
			}
			starts = getAdjacent(ctx, curSide)
		}

		// We consider any tile that can place `pips` in one of the `starts`.
		for _, t := range p.GetHand() {
			// Unfortunate duplication.
			if t.GetA() == pips {
				for _, a := range starts {
					bs := getAdjacent(ctx, a)
					for _, b := range bs {
						nextPlacement := &tpb.Placement{
							Tile: t,
							A:    a,
							B:    b,
						}
						if !mask.getp(nextPlacement) {
							moves = append(moves, nextPlacement)
						}
					}
				}
			}
			if t.GetB() == pips {
				for _, b := range starts {
					as := getAdjacent(ctx, b)
					for _, a := range as {
						nextPlacement := &tpb.Placement{
							Tile: t,
							A:    a,
							B:    b,
						}
						if !mask.getp(nextPlacement) {
							moves = append(moves, nextPlacement)
						}
					}
				}
			}
		}
	}

	// TODO: disallow moves that "crowd" the round leader.

	return moves, nil
}

type mask struct {
	v []bool
	w int32
	h int32
}

func maskForBoard(ctx context.Context, b *tpb.Board) mask {
	m := mask{
		v: make([]bool, b.GetWidth()*b.GetHeight()),
		w: b.GetWidth(),
	}
	allLines := append([]*tpb.Line{}, b.PlayerLines...)
	allLines = append(allLines, b.FreeLines...)
	for _, l := range allLines {
		for _, placement := range l.GetPlacements() {
			m.setp(placement)
		}
	}
	return m
}

func (m mask) setp(placement *tpb.Placement) {
	m.setc(placement.GetA())
	m.setc(placement.GetB())
}

func (m mask) setc(coord *tpb.Coord) {
	if coord.GetX() < 0 || coord.GetY() < 0 {
		return
	}
	if coord.GetX() >= m.w || coord.GetY() >= m.h {
		return
	}
	m.v[coord.GetX()+coord.GetY()*m.w] = true
}

func (m mask) getp(placement *tpb.Placement) bool {
	return m.getc(placement.GetA()) || m.getc(placement.GetB())
}

func (m mask) getc(coord *tpb.Coord) bool {
	if coord.GetX() < 0 || coord.GetY() < 0 {
		return true
	}
	if coord.GetX() >= m.w || coord.GetY() >= m.h {
		return true
	}
	return m.v[coord.GetX()+coord.GetY()*m.w]
}

func getAdjacent(ctx context.Context, coord *tpb.Coord) []*tpb.Coord {
	return []*tpb.Coord{
		&tpb.Coord{
			X: coord.X,
			Y: coord.Y + 1,
		},
		&tpb.Coord{
			X: coord.X,
			Y: coord.Y - 1,
		},
		&tpb.Coord{
			X: coord.X + 1,
			Y: coord.Y,
		},
		&tpb.Coord{
			X: coord.X - 1,
			Y: coord.Y,
		},
	}
}
