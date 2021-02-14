package tiles

import (
	"testing"

	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

func TestMask(t *testing.T) {
	for _, tc := range []struct {
		label      string
		width      int32
		height     int32
		placements []*tpb.Placement
		wantHits   []*tpb.Placement
		wantMisses []*tpb.Placement
	}{{
		label:  "just leader",
		width:  11,
		height: 10,
		placements: []*tpb.Placement{{
			A: &tpb.Coord{X: 4, Y: 5},
			B: &tpb.Coord{X: 5, Y: 5},
		}},
		wantHits: []*tpb.Placement{{
			A: &tpb.Coord{X: 3, Y: 5},
			B: &tpb.Coord{X: 4, Y: 5},
		}},
		wantMisses: []*tpb.Placement{{
			A: &tpb.Coord{X: 3, Y: 5},
			B: &tpb.Coord{X: 2, Y: 5},
		}},
	}} {
		t.Run(tc.label, func(t *testing.T) {
			m := mask{
				v: make([]bool, tc.width*tc.height),
				w: tc.width,
				h: tc.height,
			}
			for _, p := range tc.placements {
				m.setp(p)
			}
			for _, p := range tc.wantHits {
				if !m.getp(p) {
					t.Errorf("with %q; got miss, want hit", p)
				}
			}
			for _, p := range tc.wantMisses {
				if m.getp(p) {
					t.Errorf("with %q; got hit, want miss", p)
				}
			}
		})
	}
}
