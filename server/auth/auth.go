package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var playerIDKey struct{}

func AccessFilter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}
	tokens := md.Get("access_token")
	if len(tokens) != 1 || tokens[0] == "" {
		return nil, status.Errorf(codes.PermissionDenied, "bad access_token in metadata: %q", tokens)
	}
	accessToken := tokens[0]

	tokens = md.Get("player_id")
	if len(tokens) != 1 || tokens[0] == "" {
		return nil, status.Errorf(codes.PermissionDenied, "bad player_id in metadata: %q", tokens)
	}
	playerID := tokens[0]
	if playerID != accessToken {
		return nil, status.Error(codes.PermissionDenied, "access token invalid for player")
	}

	// Theoretically here we'd look up the player ID that corresponds to the token.
	// For now, the access_token is the player ID.
	ctx = context.WithValue(ctx, playerIDKey, accessToken)

	resp, err := handler(ctx, req)
	return resp, err
}

func PlayerIDFromContext(ctx context.Context) (string, bool) {
	pid := ctx.Value(playerIDKey)
	if pid == nil {
		return "", false
	}
	playerID, ok := pid.(string)
	return playerID, ok
}
