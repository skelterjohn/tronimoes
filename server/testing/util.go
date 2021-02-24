package testing

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"

	"github.com/skelterjohn/tronimoes/server"
	"github.com/skelterjohn/tronimoes/server/auth"
	spb "github.com/skelterjohn/tronimoes/server/proto"
)

func playMovesUntilDone(t *testing.T, ctx context.Context, c spb.TronimoesClient, playerID, gameID string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"access_token": playerID,
		"player_id":    playerID,
	}))

	for {
		// Wait until this player's turn.
		b, err := c.GetBoard(ctx, &spb.GetBoardRequest{
			GameId: gameID,
		})

		if b.GetDone() {
			return
		}

		if err != nil {
			t.Fatalf("Error getting board for %s: %v", playerID, err)
			return
		}

		if b.GetNextPlayerId() == "" {
			t.Fatal("It was no one's turn")
			return
		}

		if b.GetNextPlayerId() != playerID {
			// Since all players will be automated, we can have a short wait time to avoid a long test.
			time.Sleep(10 * time.Millisecond)
			continue
		}

		resp, err := c.GetMoves(ctx, &spb.GetMovesRequest{
			GameId: gameID,
		})
		if err != nil {
			t.Fatalf("Could not get moves: %v", err)
		}

		placements := resp.GetPlacements()
		if len(placements) == 0 {
			return
		}

		t.Logf("%s has %d moves", playerID, len(placements))

		placement := placements[rand.Intn(len(placements))]

		if _, err := c.LayTile(ctx, &spb.LayTileRequest{
			GameId:    gameID,
			Placement: placement,
		}); err != nil {
			t.Fatalf("Error laying tile: %v", err)
		}
		fmt.Printf("%s played %q\n", playerID, placement)
	}
}

func gameForPlayers(t *testing.T, ctx context.Context, c spb.TronimoesClient, playerIDs []string) string {
	t.Helper()

	games := make([]*spb.Game, len(playerIDs))

	wg := &sync.WaitGroup{}
	for i, pid := range playerIDs {
		wg.Add(1)
		go func(ctx context.Context, i int, pid string) {
			defer wg.Done()
			g, err := createGameAndWait(t, ctx, c, pid, &spb.CreateGameRequest{
				Discoverable: false,
				Private:      false,
				MinPlayers:   0,
				MaxPlayers:   0,
				BoardShape:   spb.CreateGameRequest_standard_31_by_30,
			})
			if err != nil {
				t.Fatalf("Could not create game for %s: %v", pid, err)
			}
			games[i] = g
		}(ctx, i, pid)
	}
	wg.Wait()

	for i := range games {
		if games[i].GetGameId() != games[0].GetGameId() {
			t.Fatalf("Got different games for different players, %s != %s", games[i].GetGameId(), games[0].GetGameId())
		}
	}

	return games[0].GetGameId()
}

func createGameAndWait(t *testing.T, ctx context.Context, c spb.TronimoesClient, playerID string, req *spb.CreateGameRequest) (*spb.Game, error) {
	t.Helper()

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"access_token": playerID,
		"player_id":    playerID,
	}))

	resp, err := c.CreateGame(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not create game: %v", err)
	}

	for !resp.GetDone() {
		time.Sleep(1 * time.Second)
		resp, err = c.GetOperation(ctx, &spb.GetOperationRequest{
			OperationId: resp.GetOperationId(),
		})
		if err != nil {
			return nil, fmt.Errorf("error getting operation: %v", err)
		}
	}

	if resp.GetStatus() != spb.Operation_SUCCESS {
		return nil, fmt.Errorf("operation not SUCCESS, got %q instead", resp.GetStatus())
	}

	g := &spb.Game{}
	if resp.GetPayload().GetTypeUrl() != "skelterjohn.tronimoes.Game" {
		return nil, fmt.Errorf("unexpected operation payload type %q", resp.GetPayload().GetTypeUrl())
	}
	if err := proto.Unmarshal(resp.GetPayload().GetValue(), g); err != nil {
		return nil, fmt.Errorf("could not unmarshal operation payload: %v", err)
	}

	return g, nil
}

func createBufferedServer(t *testing.T, ctx context.Context) (spb.TronimoesClient, func()) {
	t.Helper()

	l := bufconn.Listen(10 * 1024)

	operations := &server.InMemoryOperations{}
	games := &server.InMemoryGames{}
	queue := &server.InMemoryQueue{
		Games:      games,
		Operations: operations,
	}

	tronimoes := &server.Tronimoes{
		Operations: operations,
		Games:      games,
		Queue:      queue,
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(auth.AccessFilter))
	spb.RegisterTronimoesServer(s, tronimoes)
	reflection.Register(s)

	go func() {
		t.Helper()
		if err := s.Serve(l); err != nil {
			t.Errorf("Error serving: %v", err)
		}
	}()

	dial := func(context.Context, string) (net.Conn, error) {
		return l.Dial()
	}

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dial), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not dial buffered server: %v", err)
	}
	return spb.NewTronimoesClient(conn), func() {
		if err := conn.Close(); err != nil {
			t.Errorf("Error closing conn: %v", err)
		}
		if err := l.Close(); err != nil {
			t.Errorf("Error closing listener: %v", err)
		}
		s.Stop()
	}
}
