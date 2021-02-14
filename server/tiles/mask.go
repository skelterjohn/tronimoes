package tiles

import (
	"context"

	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

type mask struct {
	v []bool
	w int32
	h int32
}

func maskForBoard(ctx context.Context, b *tpb.Board) mask {
	m := mask{
		v: make([]bool, b.GetWidth()*b.GetHeight()),
		w: b.GetWidth(),
		h: b.GetHeight(),
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
