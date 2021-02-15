package main

import (
	"context"
	"crypto/tls"
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
	useTLS = flag.Bool("tls", true, "Use TLS")
)

func main() {
	flag.Parse()
	address := flag.Arg(0)

	opts := []grpc.DialOption{}

	if *useTLS {
		config := &tls.Config{}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
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
