package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/skelterjohn/tronimoes/server"
	"github.com/skelterjohn/tronimoes/server/auth"
)

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := auth.AccessFilter(ctx, req, info, handler)
	log.Printf("RPC: %s=%v latency=%v\n", info.FullMethod, status.Code(err), time.Since(start))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
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
		grpc.UnaryInterceptor(interceptor),
	)

	if err := server.Serve(ctx, port, s); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
