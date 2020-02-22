package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

const (
	defaultName = "world"
)

var (
	tls = flag.Bool("tls", true, "Use TLS")
)

func main() {
	flag.Parse()
	address := flag.Arg(0)

	creds, err := credentials.NewClientTLSFromFile("/etc/ssl/certs/ca-certificates.crt", "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	_ = creds

	opt := grpc.WithTransportCredentials(creds)
	if !*tls {
		opt = grpc.WithInsecure()
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opt)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := tpb.NewTronimoesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateGame(ctx, &tpb.CreateGameRequest{})
	if err != nil {
		log.Fatalf("could not create game: %v", err)
	}
	log.Printf("Response: %q", r.GetOperationId())
}
