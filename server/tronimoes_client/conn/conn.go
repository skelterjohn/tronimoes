package conn

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

func GetClient(ctx context.Context, address string, useTLS bool) (tpb.TronimoesClient, error) {
	if useTLS && !strings.Contains(address, ":") {
		address += ":443"
	}

	opts := []grpc.DialOption{}

	if useTLS {
		config := &tls.Config{}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	return tpb.NewTronimoesClient(conn), nil
}
