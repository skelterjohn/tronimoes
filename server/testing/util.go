package testing

import (
	"context"
	"fmt"
	"net"
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
