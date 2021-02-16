package testing

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/skelterjohn/tronimoes/server"
	spb "github.com/skelterjohn/tronimoes/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

var (
	serverAddress = getServerAddress()
	useTLS        = getUseTLS()
)

func getServerAddress() string {
	addr := os.Getenv("TRONIMOES_TESTING_SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:8082"
	}
	return addr
}

func getUseTLS() bool {
	tstr := os.Getenv("TRONIMOES_TESTING_USE_TLS")
	if tstr == "" {
		return false
	}
	return tstr == "1"
}

func createGameAndWait(t *testing.T, ctx context.Context, c spb.TronimoesClient, req *spb.CreateGameRequest) (*spb.Game, error) {
	t.Helper()

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

	tronimoes := &server.Tronimoes{
		Operations: operations,
		GameQueue: &server.InMemoryQueue{
			Games:      &server.InMemoryGames{},
			Operations: operations,
		},
	}

	s := grpc.NewServer()
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

func TestCreate(t *testing.T) {
	ctx := context.Background()
	c, close := createBufferedServer(t, ctx)
	defer close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var g1, g2 *spb.Game
	var err1, err2 error
	go func() {
		defer wg.Done()
		g1, err1 = createGameAndWait(t, ctx, c, &spb.CreateGameRequest{
			Discoverable: false,
			Private:      false,
			MinPlayers:   0,
			MaxPlayers:   0,
			PlayerId:     "jt",
		})
	}()
	go func() {
		defer wg.Done()
		g2, err2 = createGameAndWait(t, ctx, c, &spb.CreateGameRequest{
			Discoverable: false,
			Private:      false,
			MinPlayers:   0,
			MaxPlayers:   0,
			PlayerId:     "stef",
		})
	}()
	wg.Wait()
	if err1 != nil {
		t.Fatalf("Could not create game 1: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("Could not create game 2: %v", err2)
	}

	if g1.GetGameId() != g2.GetGameId() {
		t.Errorf("Game IDs did not match, %q != %q", g1.GetGameId(), g2.GetGameId())
	}
}
