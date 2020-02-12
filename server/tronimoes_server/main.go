package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/skelterjohn/tronimoes/server"
	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

func main() {
	// PORT is being set by the Cloud Run environment
	port := os.Getenv("PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	tpb.RegisterGameServer(s, &server.Game{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
