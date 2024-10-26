package pq

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	tpb "github.com/skelterjohn/tronimoes/server/tiles/proto"
	"github.com/skelterjohn/tronimoes/server/util"
)

type PQGames struct {
	DB *sql.DB
}

func (g *PQGames) WriteGame(ctx context.Context, gm *spb.Game) error {
	playersData := make([][]byte, len(gm.GetPlayers()))
	for i, p := range gm.GetPlayers() {
		var err error
		playersData[i], err = proto.Marshal(p)
		if err != nil {
			return util.Annotate(err, "could not marshal player")
		}
	}

	if _, err := g.DB.Exec(`
		INSERT INTO games (game_id, players, status, round_leaders, board_shape)
			VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (game_id)
		DO
			UPDATE games SET (board) = ($2, $3, $4, $5)`,
		gm.GetGameId(), playersData, gm.GetStatus().String(), gm.GetRounderLeaders(), gm.GetBoardShape().String()); err != nil {
		return util.Annotate(err, "could not write game")
	}

	return nil
}

func (g *PQGames) ReadGame(ctx context.Context, id string) (*spb.Game, error) {
	rows, err := g.DB.Query(`
		SELECT
			players,
			status,
			round_leaders,
			board_shape
		FROM games
		WHERE game_id = $1`,
		id)
	if err != nil {
		return nil, util.Annotate(err, "could not read board")
	}

	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "board not found")
	}

	var playersData [][]byte
	var statusString, boardShapeString string
	tg := &spb.Game{}

	if err := rows.Scan(&playersData, &statusString, &tg.RoundLeaders, &boardShapeString); err != nil {
		return nil, util.Annotate(err, "could not scan row into game")
	}

	b := &tpb.Board{}
	err = proto.Unmarshal(boardData, b)
	if err != nil {
		return nil, util.Annotate(err, "could not unmarshal board data")
	}

	return b, nil
}

func (g *PQGames) WriteBoard(ctx context.Context, id string, b *tpb.Board) error {
	boardData, err := proto.Marshal(b)
	if err != nil {
		return util.Annotate(err, "could not marshal board")
	}

	if _, err := g.DB.Exec(`
		INSERT INTO boards (game_id, board)
			VALUES ($1, $2)
		ON CONFLICT (game_id)
		DO
			UPDATE boards SET (board) = ($2)`,
		id, boardData); err != nil {
		return util.Annotate(err, "could not write board")
	}

	return nil
}

func (g *PQGames) ReadBoard(ctx context.Context, id string) (*tpb.Board, error) {
	rows, err := g.DB.Query(`
		SELECT
			board
		FROM boards
		WHERE game_id = $1`,
		id)
	if err != nil {
		return nil, util.Annotate(err, "could not read board")
	}

	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "board not found")
	}

	var boardData []byte

	if err := rows.Scan(&boardData); err != nil {
		return nil, util.Annotate(err, "could not scan row into boardData")
	}

	b := &tpb.Board{}
	err = proto.Unmarshal(boardData, b)
	if err != nil {
		return nil, util.Annotate(err, "could not unmarshal board data")
	}

	return b, nil
}
