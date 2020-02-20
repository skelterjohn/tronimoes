package server

import (
	"context"
	"log"
	"os"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Game struct {
}

func (g *Game) Hello(ctx context.Context, req *tpb.HelloRequest) (*tpb.HelloResponse, error) {
	log.Printf("Returning response for %q from v%s", req.Message, os.Getenv("SHORT_SHA"))
	return &tpb.HelloResponse{
		Message:  "Echo: " + req.Message,
		Revision: "v" + os.Getenv("SHORT_SHA"),
	}, nil
}
