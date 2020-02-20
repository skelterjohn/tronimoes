package server

import (
	"context"
	"os"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Game struct {
}

func (g *Game) Hello(ctx context.Context, req *tpb.HelloRequest) (*tpb.HelloResponse, error) {
	return &tpb.HelloResponse{
		Message:  "Echo: " + req.Message,
		Revision: os.Getenv("SHORT_SHA"),
	}, nil
}
