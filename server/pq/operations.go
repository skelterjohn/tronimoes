package pq

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PQOperations struct {
	DB *sql.DB
}

/*

message Operation {
	string operation_id = 1;

	bool done = 2;

	google.protobuf.Any payload = 3;

	enum Status {
		UNKNOWN = 0;
		SUCCESS = 1;
		FAILURE = 2;
	}
	Status status = 4;
}
*/

func (o *PQOperations) WriteOperation(ctx context.Context, op *spb.Operation) error {
	_, err := o.DB.Exec(`
		UPDATE operations
		SET (done, payload, status) = ($1, $2, $3)
		WHERE operation_id=$4`,
		op.GetDone(), op.GetPayload(), op.GetStatus(), op.GetOperationId())

	if err != nil {
		return util.Annotate(err, "could not write operation")
	}

	return nil
}

func (o *PQOperations) ReadOperation(ctx context.Context, id string) (*spb.Operation, error) {
	rows, err := o.DB.Query(`
		SELECT (done, payload, status)
		FROM operations
		WHERE operation_id = $1`,
		id)
	if err != nil {
		return nil, util.Annotate(err, "could not read operation")
	}

	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "operation not found")
	}

	op := &spb.Operation{
		OperationId: id,
	}

	if err := rows.Scan(&op.Done, &op.Payload, &op.Status); err != nil {
		return nil, util.Annotate(err, "could not scan row into operation")
	}

	return op, nil
}

func (o *PQOperations) NewOperation(ctx context.Context) (*spb.Operation, error) {
	op := &spb.Operation{
		OperationId: uuid.New().String(),
	}
	_, err := o.DB.Exec(`
		INSERT INTO operations
		(operation_id, done, payload, status) VALUES ($1, $2, $3, $4)`,
		op.GetOperationId(), op.GetDone(), op.GetPayload(), op.GetStatus())

	if err != nil {
		return nil, util.Annotate(err, "could not insert operation")
	}
	return op, nil
}
