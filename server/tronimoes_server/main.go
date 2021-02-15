package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/skelterjohn/tronimoes/server"
	spb "github.com/skelterjohn/tronimoes/server/proto"
)

func RPCSummary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("RPC: %s=%v latency=%v", info.FullMethod, status.Code(err), time.Since(start))
	return resp, err
}

func main() {
	// PORT is being set by the Cloud Run environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(RPCSummary),
	)

	operations := &server.InMemoryOperations{}

	tronimoes := &server.Tronimoes{
		Operations: operations,
		GameQueue: &server.InMemoryQueue{
			Games:      &server.InMemoryGames{},
			Operations: operations,
		},
	}

	spb.RegisterTronimoesServer(s, tronimoes)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
