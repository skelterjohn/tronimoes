package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

const (
	defaultName = "world"
)

func main() {
	address := os.Getenv("ADDR")

	creds, err := credentials.NewClientTLSFromFile("/etc/ssl/certs/ca-certificates.crt", "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := tpb.NewGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &tpb.HelloRequest{Message: "test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
