package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/skelterjohn/tronimoes/server"
)

func RPCSummary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("RPC: %s=%v latency=%v", info.FullMethod, status.Code(err), time.Since(start))
	return resp, err
}

func main() {
	ctx := context.Background()

	// PORT is being set by the Cloud Run environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(RPCSummary),
	)

	if err := server.Serve(ctx, port, s); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
