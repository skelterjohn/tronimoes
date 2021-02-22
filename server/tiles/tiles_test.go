package tiles

import (
	"bytes"
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

func TestSetupBoard(t *testing.T) {
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

func tab(a, b int32) *tpb.Tile {
	return &tpb.Tile{
		A: a,
		B: b,
	}
}
func xy(x, y int32) *tpb.Coord {
	return &tpb.Coord{
		X: x,
		Y: y,
	}
}

func sortedBoard() *tpb.Board {
	return &tpb.Board{
		Players: []*tpb.Player{{
			PlayerId: "jt",
		}, {
			PlayerId: "stef",
		}},
		Bag:    sortedBag(12),
		Width:  11,
		Height: 10,
	}
}

func TestEdgeRace(t *testing.T) {
	Debug = true
	ctx := context.Background()

	board := sortedBoard()

	board, err := SetupBoard(ctx, board, 100)
	if err != nil {
		t.Fatal(err)
	}
	board.GetPlayers()[0].Hand = []*tpb.Tile{
		tab(12, 12), tab(12, 11), tab(11, 10), tab(10, 9), tab(9, 8),
	}
	board.GetPlayers()[1].Hand = []*tpb.Tile{
		tab(12, 10), tab(10, 8), tab(8, 6), tab(6, 4), tab(4, 2),
	}

	for i, l := range board.GetPlayerLines() {
		if len(l.GetPlacements()) != 1 {
			t.Fatalf("Bad initial line for %s", board.GetPlayers()[i].GetPlayerId())
		}
		if !proto.Equal(l.GetPlacements()[0].GetTile(), tab(12, 12)) {
			t.Fatalf("Wrong round leader %q for %s", l.GetPlacements()[0], board.GetPlayers()[i].GetPlayerId())
		}
	}

	placements := []*tpb.Placement{{
		Tile: tab(12, 10),
		A:    xy(6, 5),
		B:    xy(7, 5),
	}, {
		Tile: tab(12, 11),
		A:    xy(3, 5),
		B:    xy(2, 5),
	}, {
		Tile: tab(10, 8),
		A:    xy(8, 5),
		B:    xy(9, 5),
	}, {
		Tile: tab(11, 10),
		A:    xy(1, 5),
		B:    xy(0, 5),
	}}

	// Play a bunch of legal moves.

	for _, placement := range placements {
		if err := LayTile(ctx, board, placement); err != nil {
			t.Fatalf("Error placing %q: %v", placement, err)
		}
	}

	// Try to play a move using a tile not in the hand.
	if err := LayTile(ctx, board, &tpb.Placement{
		Tile: tab(8, 1),
		A:    xy(9, 4),
		B:    xy(9, 3),
	}); err == nil {
		t.Fatal("able to lay tile not in hand")
	}

	// Try to play a tile somewhere disconnected.
	if err := LayTile(ctx, board, &tpb.Placement{
		Tile: tab(8, 6),
		A:    xy(1, 2),
		B:    xy(1, 3),
	}); err == nil {
		t.Fatal("able to lay disconnected tile")
	}
}

func TestLegalMoves(t *testing.T) {
	Debug = true
	ctx := context.Background()
	newBoard, err := SetupBoard(ctx, sortedBoard(), 13)
	if err != nil {
		t.Fatalf("Problem making new board: %v", err)
	}

	for _, tc := range []struct {
		label string
		board *tpb.Board
		want  []*tpb.Placement
	}{{
		label: "firstmove",
		board: newBoard,
		want: []*tpb.Placement{{
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 4, Y: 7},
			B:    &tpb.Coord{X: 4, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 5, Y: 6},
			B:    &tpb.Coord{X: 4, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 3, Y: 6},
			B:    &tpb.Coord{X: 4, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 4, Y: 3},
			B:    &tpb.Coord{X: 4, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 5, Y: 4},
			B:    &tpb.Coord{X: 4, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 3, Y: 4},
			B:    &tpb.Coord{X: 4, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 3, Y: 6},
			B:    &tpb.Coord{X: 3, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 3, Y: 4},
			B:    &tpb.Coord{X: 3, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 2, Y: 5},
			B:    &tpb.Coord{X: 3, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 5, Y: 7},
			B:    &tpb.Coord{X: 5, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 6, Y: 6},
			B:    &tpb.Coord{X: 5, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 4, Y: 6},
			B:    &tpb.Coord{X: 5, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 5, Y: 3},
			B:    &tpb.Coord{X: 5, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 6, Y: 4},
			B:    &tpb.Coord{X: 5, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 4, Y: 4},
			B:    &tpb.Coord{X: 5, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 6, Y: 6},
			B:    &tpb.Coord{X: 6, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 6, Y: 4},
			B:    &tpb.Coord{X: 6, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 8, B: 12},
			A:    &tpb.Coord{X: 7, Y: 5},
			B:    &tpb.Coord{X: 6, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 4, Y: 7},
			B:    &tpb.Coord{X: 4, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 5, Y: 6},
			B:    &tpb.Coord{X: 4, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 3, Y: 6},
			B:    &tpb.Coord{X: 4, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 4, Y: 3},
			B:    &tpb.Coord{X: 4, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 5, Y: 4},
			B:    &tpb.Coord{X: 4, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 3, Y: 4},
			B:    &tpb.Coord{X: 4, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 3, Y: 6},
			B:    &tpb.Coord{X: 3, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 3, Y: 4},
			B:    &tpb.Coord{X: 3, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 2, Y: 5},
			B:    &tpb.Coord{X: 3, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 5, Y: 7},
			B:    &tpb.Coord{X: 5, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 6, Y: 6},
			B:    &tpb.Coord{X: 5, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 4, Y: 6},
			B:    &tpb.Coord{X: 5, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 5, Y: 3},
			B:    &tpb.Coord{X: 5, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 6, Y: 4},
			B:    &tpb.Coord{X: 5, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 4, Y: 4},
			B:    &tpb.Coord{X: 5, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 6, Y: 6},
			B:    &tpb.Coord{X: 6, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 6, Y: 4},
			B:    &tpb.Coord{X: 6, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 7, B: 12},
			A:    &tpb.Coord{X: 7, Y: 5},
			B:    &tpb.Coord{X: 6, Y: 5},
		}},
	}, {
		label: "next move",
		board: &tpb.Board{
			Bag: newBoard.GetBag(),
			Players: []*tpb.Player{{
				PlayerId: "jt",
				Hand: []*tpb.Tile{{
					A: 10,
					B: 11,
				}},
			}},
			NextPlayerId: "jt",
			Width:        11,
			Height:       10,
			PlayerLines: []*tpb.Line{{
				PlayerId: "jt",
				Placements: []*tpb.Placement{{
					Tile: &tpb.Tile{A: 12, B: 12},
					A:    &tpb.Coord{X: 4, Y: 5},
					B:    &tpb.Coord{X: 5, Y: 5},
				}, {
					Tile: &tpb.Tile{A: 11, B: 12},
					A:    &tpb.Coord{X: 2, Y: 5},
					B:    &tpb.Coord{X: 3, Y: 5},
				}},
			}},
		},
		want: []*tpb.Placement{{
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 2, Y: 7},
			B:    &tpb.Coord{X: 2, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 3, Y: 6},
			B:    &tpb.Coord{X: 2, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 1, Y: 6},
			B:    &tpb.Coord{X: 2, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 2, Y: 3},
			B:    &tpb.Coord{X: 2, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 3, Y: 4},
			B:    &tpb.Coord{X: 2, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 1, Y: 4},
			B:    &tpb.Coord{X: 2, Y: 4},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 1, Y: 6},
			B:    &tpb.Coord{X: 1, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 1, Y: 4},
			B:    &tpb.Coord{X: 1, Y: 5},
		}, {
			Tile: &tpb.Tile{A: 10, B: 11},
			A:    &tpb.Coord{X: 0, Y: 5},
			B:    &tpb.Coord{X: 1, Y: 5},
		}},
	}, {
		label: "next move adjacent",
		board: &tpb.Board{
			Bag: newBoard.GetBag(),
			Players: []*tpb.Player{{
				PlayerId: "jt",
				Hand: []*tpb.Tile{{
					A: 10,
					B: 8,
				}},
			}, {
				PlayerId: "stef",
				Hand: []*tpb.Tile{{
					A: 11,
					B: 9,
				}},
			}},
			NextPlayerId: "stef",
			Width:        11,
			Height:       10,
			PlayerLines: []*tpb.Line{{
				PlayerId: "jt",
				Placements: []*tpb.Placement{{
					Tile: &tpb.Tile{A: 12, B: 12},
					A:    &tpb.Coord{X: 4, Y: 5},
					B:    &tpb.Coord{X: 5, Y: 5},
				}, {
					Tile: &tpb.Tile{A: 11, B: 12},
					A:    &tpb.Coord{X: 4, Y: 6},
					B:    &tpb.Coord{X: 4, Y: 7},
				}, {
					Tile: &tpb.Tile{A: 11, B: 10},
					A:    &tpb.Coord{X: 4, Y: 8},
					B:    &tpb.Coord{X: 4, Y: 9},
				}},
			}, {
				PlayerId: "stef",
				Placements: []*tpb.Placement{{
					Tile: &tpb.Tile{A: 12, B: 12},
					A:    &tpb.Coord{X: 4, Y: 5},
					B:    &tpb.Coord{X: 5, Y: 5},
				}, {
					Tile: &tpb.Tile{A: 12, B: 11},
					A:    &tpb.Coord{X: 5, Y: 6},
					B:    &tpb.Coord{X: 5, Y: 7},
				}},
			}},
		},
		want: []*tpb.Placement{{
			Tile: &tpb.Tile{A: 11, B: 9},
			A:    &tpb.Coord{X: 5, Y: 8},
			B:    &tpb.Coord{X: 5, Y: 9},
		}, {
			Tile: &tpb.Tile{A: 11, B: 9},
			A:    &tpb.Coord{X: 5, Y: 8},
			B:    &tpb.Coord{X: 6, Y: 8},
		}, {
			Tile: &tpb.Tile{A: 11, B: 9},
			A:    &tpb.Coord{X: 6, Y: 7},
			B:    &tpb.Coord{X: 6, Y: 8},
		}, {
			Tile: &tpb.Tile{A: 11, B: 9},
			A:    &tpb.Coord{X: 6, Y: 7},
			B:    &tpb.Coord{X: 6, Y: 6},
		}, {
			Tile: &tpb.Tile{A: 11, B: 9},
			A:    &tpb.Coord{X: 6, Y: 7},
			B:    &tpb.Coord{X: 7, Y: 7},
		}},
	}} {
		t.Run(tc.label, func(t *testing.T) {
			got, err := LegalMoves(ctx, tc.board)
			if err != nil {
				t.Fatalf("Could not get legal moves: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Errorf("Wrong number of placements; got %d, want %d", len(got), len(tc.want))
				t.Fatalf("Bad placements;\n got: %s\nwant: %s\n", goCodeForPlacements(got), goCodeForPlacements(tc.want))
			}
			for i := range tc.want {
				if !proto.Equal(tc.want[i], got[i]) {
					t.Fatalf("Bad placements;\n got: %s\nwant: %s\n", goCodeForPlacements(got), goCodeForPlacements(tc.want))
				}
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

func goCodeForPlacements(placements []*tpb.Placement) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "[]*tpb.Placement{{\n")
	for i, p := range placements {
		if i != 0 {
			fmt.Fprint(buf, "},{\n")
		}
		fmt.Fprintf(buf, "\tTile: &tpb.Tile{A: %d, B: %d},\n", p.GetTile().GetA(), p.GetTile().GetB())
		fmt.Fprintf(buf, "\tA: &tpb.Coord{X: %d, Y: %d},\n", p.GetA().GetX(), p.GetA().GetY())
		fmt.Fprintf(buf, "\tB: &tpb.Coord{X: %d, Y: %d},\n", p.GetB().GetX(), p.GetB().GetY())
	}
	fmt.Fprintf(buf, "}}\n")
	return buf.String()
}
