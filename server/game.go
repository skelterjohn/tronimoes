package server

import (
	"context"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type Game struct {
}

func (g *Game) Hello(context.Context, *tpb.HelloRequest) (*tpb.HelloResponse, error) {
	return nil, nil
}
