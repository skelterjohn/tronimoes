package conn

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	spb "github.com/skelterjohn/tronimoes/server/proto"
)

func GetClient(ctx context.Context, address string, useTLS bool) (spb.TronimoesClient, error) {
	if !strings.Contains(address, ":") {
		if useTLS {
			address += ":443"
		} else {
			address += ":8082"
		}
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

	return spb.NewTronimoesClient(conn), nil
}
