package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/skelterjohn/tronimoes/server/auth"
	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/tiles"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
)

func annotatef(err error, format string, items ...interface{}) error {
	upstream := err.Error()
	if s, ok := status.FromError(err); ok {
		upstream = s.Message()
	}
	msg := fmt.Sprintf(format, items...) + ": " + upstream
	return status.Error(status.Code(err), msg)
}

type Operations interface {
	WriteOperation(ctx context.Context, op *spb.Operation) error
	ReadOperation(ctx context.Context, id string) (*spb.Operation, error)
	NewOperation(ctx context.Context) (*spb.Operation, error)
}

type Queue interface {
	AddPlayer(ctx context.Context, playerID string, req *spb.CreateGameRequest, operationID string) error
	MakeNextGame(ctx context.Context) error
}

type Tronimoes struct {
	Operations Operations
	Queue      Queue
	Games      Games
	Rounds     *Rounds
}

func (t *Tronimoes) CreateAccessToken(ctx context.Context, req *spb.CreateAccessTokenRequest) (*spb.AccessTokenResponse, error) {
	exp, err := ptypes.TimestampProto(time.Now().Add(24 * time.Hour))
	if err != nil {
		exp = ptypes.TimestampNow()
	}
	return &spb.AccessTokenResponse{
		AccessToken: req.GetPlayerId(),
		Expiry:      exp,
	}, nil
}

func (t *Tronimoes) CreateGame(ctx context.Context, req *spb.CreateGameRequest) (*spb.Operation, error) {
	playerID, ok := auth.PlayerIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unknown player ID")
	}
	op, err := t.Operations.NewOperation(ctx)
	if err != nil {
		return nil, annotatef(err, "could not create operation")
	}
	if err := t.Queue.AddPlayer(ctx, playerID, req, op.GetOperationId()); err != nil {
		return nil, annotatef(err, "could not create queue player")
	}

	if err := t.Queue.MakeNextGame(ctx); err != nil && status.Code(err) != codes.NotFound {
		log.Printf(annotatef(err, "could not find the next game").Error())
	}

	return t.Operations.ReadOperation(ctx, op.GetOperationId())
}

func (t *Tronimoes) GetOperation(ctx context.Context, req *spb.GetOperationRequest) (*spb.Operation, error) {
	return t.Operations.ReadOperation(ctx, req.GetOperationId())
}

func (t *Tronimoes) checkPlayerInGame(ctx context.Context, playerID string, g *spb.Game) error {

	foundPlayer := false
	for _, p := range g.GetPlayers() {
		if p.GetPlayerId() == playerID {
			foundPlayer = true
		}
	}
	if foundPlayer {
		return nil
	}
	return status.Errorf(codes.PermissionDenied, "%s is not a player in game %s", playerID, g.GetGameId())
}

func (t *Tronimoes) GetGame(ctx context.Context, req *spb.GetGameRequest) (*spb.Game, error) {
	playerID, ok := auth.PlayerIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unknown player ID")
	}
	g, err := t.Games.ReadGame(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get game")
	}
	if err := t.checkPlayerInGame(ctx, playerID, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (t *Tronimoes) GetBoard(ctx context.Context, req *spb.GetBoardRequest) (*tpb.Board, error) {
	playerID, ok := auth.PlayerIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unknown player ID")
	}

	g, err := t.Games.ReadGame(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get game")
	}
	if err := t.checkPlayerInGame(ctx, playerID, g); err != nil {
		return nil, err
	}

	b, err := t.Games.ReadBoard(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get board")
	}

	// Nil out all tiles in the bag and other players' hands.
	for i := range b.GetBag() {
		b.Bag[i].A = -1
		b.Bag[i].B = -1
	}
	for _, p := range b.GetPlayers() {
		if p.GetPlayerId() == playerID {
			continue
		}
		for i := range p.GetHand() {
			p.Hand[i].A = -1
			p.Hand[i].B = -1
		}
	}

	return b, nil
}

func (t *Tronimoes) GetMoves(ctx context.Context, req *spb.GetMovesRequest) (*spb.GetMovesResponse, error) {
	playerID, ok := auth.PlayerIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unknown player ID")
	}

	g, err := t.Games.ReadGame(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get game")
	}
	if err := t.checkPlayerInGame(ctx, playerID, g); err != nil {
		return nil, err
	}

	b, err := t.Games.ReadBoard(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get board")
	}

	if b.GetNextPlayerId() != playerID {
		return nil, status.Error(codes.FailedPrecondition, "it is not your turn")
	}

	var player *tpb.Player
	for _, p := range b.GetPlayers() {
		if p.GetPlayerId() == playerID {
			player = p
		}
	}
	if player == nil {
		return nil, status.Errorf(codes.InvalidArgument, "player %s is not in this game", playerID)
	}

	moves, err := tiles.LegalMoves(ctx, b, player)
	if err != nil {
		return nil, annotatef(err, "could not get legal moves")
	}

	return &spb.GetMovesResponse{
		Placements: moves,
	}, nil
}

func (t *Tronimoes) LayTile(ctx context.Context, req *spb.LayTileRequest) (*tpb.Board, error) {
	playerID, ok := auth.PlayerIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unknown player ID")
	}

	g, err := t.Games.ReadGame(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get game")
	}
	if err := t.checkPlayerInGame(ctx, playerID, g); err != nil {
		return nil, err
	}

	b, err := t.Games.ReadBoard(ctx, req.GetGameId())
	if err != nil {
		return nil, annotatef(err, "could not get board")
	}

	if b.GetNextPlayerId() != playerID {
		return nil, status.Error(codes.FailedPrecondition, "it is not your turn")
	}

	b, err = tiles.LayTile(ctx, b, req.GetPlacement())
	if err != nil {
		return nil, annotatef(err, "could not lay tile")
	}

	if err := t.Games.WriteBoard(ctx, req.GetGameId(), b); err != nil {
		return nil, annotatef(err, "could not write board")
	}

	return b, nil
}

func Serve(ctx context.Context, port string, s *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":"+port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	operations := &InMemoryOperations{}
	games := &InMemoryGames{}
	rounds := &Rounds{
		Games: games,
	}
	queue := &InMemoryQueue{
		Games:      games,
		Operations: operations,
		Rounds:     rounds,
	}

	tronimoes := &Tronimoes{
		Operations: operations,
		Games:      games,
		Queue:      queue,
		Rounds:     rounds,
	}

	spb.RegisterTronimoesServer(s, tronimoes)
	reflection.Register(s)

	return s.Serve(lis)
}
