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

func (o *PQOperations) WriteOperation(ctx context.Context, op *spb.Operation) error {
	_, err := o.DB.Exec(`
		UPDATE Operations
		SET (done, payload, status) = ($1, $2, $3)
		WHERE operation_id=$4`,
		op.GetDone(), op.GetPayload(), op.GetStatus().String(), op.GetOperationId())

	if err != nil {
		return util.Annotate(err, "could not write operation")
	}

	return nil
}

func (o *PQOperations) ReadOperation(ctx context.Context, id string) (*spb.Operation, error) {
	rows, err := o.DB.Query(`
		SELECT
			done,
			payload,
			status
		FROM Operations
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

	var statusString string

	if err := rows.Scan(&op.Done, &op.Payload, &statusString); err != nil {
		return nil, util.Annotate(err, "could not scan row into operation")
	}

	op.Status = spb.Operation_Status(spb.Operation_Status_value[statusString])

	return op, nil
}

func (o *PQOperations) NewOperation(ctx context.Context) (*spb.Operation, error) {
	op := &spb.Operation{
		OperationId: uuid.New().String(),
	}
	_, err := o.DB.Exec(`
		INSERT INTO Operations
		(operation_id, done, payload, status) VALUES ($1, $2, $3, $4)`,
		op.GetOperationId(), op.GetDone(), op.GetPayload(), op.GetStatus().String())

	if err != nil {
		return nil, util.Annotate(err, "could not insert operation")
	}
	return op, nil
}
