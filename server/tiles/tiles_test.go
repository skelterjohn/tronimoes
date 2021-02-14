package tiles

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

func sortedBag(upTo int32) []*tpb.Tile {
	b := []*tpb.Tile{}
	for i := upTo; i >= 0; i-- {
		for j := upTo; j >= i; j-- {
			b = append(b, &tpb.Tile{
				A: int32(i),
				B: int32(j),
			})
		}
	}
	return b
}

func TestNewBoard(t *testing.T) {
	for _, tc := range []struct {
		label      string
		board      *tpb.Board
		lastLeader int32
		expected   *tpb.Board
		err        bool
	}{{
		label: "not enough players",
		board: &tpb.Board{
			Players: []*tpb.Player{{
				PlayerId: "jt",
			}},
		},
		err: true,
	}, {
		label: "too many players",
		board: &tpb.Board{
			Players: []*tpb.Player{{
				PlayerId: "jt1",
			}, {
				PlayerId: "jt2",
			}, {
				PlayerId: "jt3",
			}, {
				PlayerId: "jt4",
			}, {
				PlayerId: "jt5",
			}, {
				PlayerId: "jt6",
			}, {
				PlayerId: "jt7",
			}},
		},
		err: true,
	}, {
		label: "two players",
		board: &tpb.Board{
			Players: []*tpb.Player{{
				PlayerId: "jt",
			}, {
				PlayerId: "stef",
			}},
			Bag:    sortedBag(12),
			Width:  11,
			Height: 10,
		},
		lastLeader: 13,
		expected: &tpb.Board{
			NextPlayerId: "stef",
			Width:        11,
			Height:       10,
			PlayerLines: []*tpb.Line{{
				PlayerId: "jt",
				Placements: []*tpb.Placement{{
					Type: tpb.Placement_PLAYER_LEADER,
					A:    &tpb.Coord{X: 4, Y: 5},
					B:    &tpb.Coord{X: 5, Y: 5},
					Tile: &tpb.Tile{
						A: 12,
						B: 12,
					},
				}},
			}, {
				PlayerId: "stef",
				Placements: []*tpb.Placement{{
					Type: tpb.Placement_PLAYER_LEADER,
					A:    &tpb.Coord{X: 4, Y: 5},
					B:    &tpb.Coord{X: 5, Y: 5},
					Tile: &tpb.Tile{
						A: 12,
						B: 12,
					},
				}},
			}},
			Players: []*tpb.Player{{
				PlayerId: "jt",
				Hand: []*tpb.Tile{{
					A: 11,
					B: 12,
				}, {
					A: 11,
					B: 11,
				}, {
					A: 10,
					B: 12,
				}, {
					A: 10,
					B: 11,
				}, {
					A: 10,
					B: 10,
				}, {
					A: 9,
					B: 12,
				}, {
					A: 9,
					B: 11,
				}, {
					A: 9,
					B: 10,
				}, {
					A: 9,
					B: 9,
				}},
			}, {
				PlayerId: "stef",
				Hand: []*tpb.Tile{{
					A: 8,
					B: 12,
				}, {
					A: 8,
					B: 11,
				}, {
					A: 8,
					B: 10,
				}, {
					A: 8,
					B: 9,
				}, {
					A: 8,
					B: 8,
				}, {
					A: 7,
					B: 12,
				}, {
					A: 7,
					B: 11,
				}, {
					A: 7,
					B: 10,
				}, {
					A: 7,
					B: 9,
				}, {
					A: 7,
					B: 8,
				}},
			}},
		},
	},
	} {
		t.Run(tc.label, func(t *testing.T) {
			ctx := context.Background()

			savedBag := tc.board.GetBag()

			result, err := SetupBoard(ctx, tc.board, tc.lastLeader)
			if tc.err {
				if err == nil {
					t.Error("Expected error, got none")
					return
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			assertBoardsEqual(t, result, tc.expected)

			takenTiles := 1
			for i := range result.GetPlayers() {
				takenTiles += len(result.GetPlayers()[i].GetHand())
				takenTiles += len(result.GetPlayerLines()[i].GetPlacements()) - 1
			}
			for _, l := range result.GetFreeLines() {
				takenTiles += len(l.GetPlacements())
			}

			expectedBag := savedBag[takenTiles:]
			if len(expectedBag) != len(result.GetBag()) {
				t.Errorf("Wrong bag size; got %d, want %d", len(result.GetBag()), len(expectedBag))
				t.Logf(" got bag:\n%q\n", result.GetBag())
				t.Logf("want bag:\n%q\n", expectedBag)
			}
		})
	}
}

func assertBoardsEqual(t *testing.T, got, want *tpb.Board) {
	t.Helper()
	if proto.Equal(got, want) {
		return
	}

	if len(got.GetPlayers()) != len(want.GetPlayers()) {
		t.Errorf("Wrong player count; got %d, want %d", len(got.GetPlayers()), len(want.GetPlayers()))
	}
	for i, gp := range got.GetPlayers() {
		if i >= len(want.GetPlayers()) {
			break
		}
		wp := want.GetPlayers()[i]
		if proto.Equal(gp, wp) {
			continue
		}
		if gp.GetName() != wp.GetName() {
			t.Errorf("Player %d wrong name; got %q, want %q", i, gp.GetName(), wp.GetName())
		}
		if gp.GetPlayerId() != wp.GetPlayerId() {
			t.Errorf("Player %d wrong player_id; got %q, want %q", i, gp.GetPlayerId(), wp.GetPlayerId())
		}
		if gp.GetChickenFooted() != wp.GetChickenFooted() {
			t.Errorf("Player %d wrong chicken_footed; got %v, want %v", i, gp.GetChickenFooted(), wp.GetChickenFooted())
		}

		// There is nothing left in the player proto, hands must be wrong.
		t.Errorf("Played %d wrong hand;\n got: %q\nwant: %q", i, gp.GetHand(), wp.GetHand())
	}

	if len(got.GetPlayerLines()) != len(want.GetPlayerLines()) {
		t.Errorf("Wrong player_line count; got %d, want %d", len(got.GetPlayerLines()), len(want.GetPlayerLines()))
	}
	for i, gl := range got.GetPlayerLines() {
		if i >= len(want.GetPlayerLines()) {
			break
		}
		wl := want.GetPlayerLines()[i]
		assertLinesEqual(t, fmt.Sprintf("player %d", i), gl, wl)
	}
	for i, gl := range got.GetFreeLines() {
		if i >= len(want.GetFreeLines()) {
			break
		}
		wl := want.GetFreeLines()[i]
		assertLinesEqual(t, fmt.Sprintf("free %d", i), gl, wl)
	}

	if got.GetNextPlayerId() != want.GetNextPlayerId() {
		t.Errorf("Wrong next_player_id; got %q, want %q", got.GetNextPlayerId(), want.GetNextPlayerId())
	}

	if got.GetWidth() != want.GetWidth() {
		t.Errorf("Wrong width; got %d, want %d", got.GetWidth(), want.GetWidth())
	}
	if got.GetHeight() != want.GetHeight() {
		t.Errorf("Wrong height; got %d, want %d", got.GetHeight(), want.GetHeight())
	}

	// We don't test the bag due to it being a pain to make test cases.
}

func assertLinesEqual(t *testing.T, which string, got, want *tpb.Line) {
	t.Helper()
	if proto.Equal(got, want) {
		return
	}

	if got.GetPlayerId() != want.GetPlayerId() {
		t.Errorf("Line %s wrong player_id; got %q, want %q", which, got.GetPlayerId(), want.GetPlayerId())
	}

	t.Errorf("Line %s wrong placements;\n got: %q\nwant: %q", which, got.GetPlacements(), want.GetPlacements())
}
