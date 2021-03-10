package pq

import (
	"context"
	"database/sql"
	"fmt"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/util"
)

type PQOperations struct {
	DB *sql.DB
}

func (o *PQOperations) WriteOperation(ctx context.Context, op *spb.Operation) error {
	rows, err := o.DB.Query(`SELECT (operation_id, done) FROM conductor.operations`)
	if err != nil {
		return util.Annotate(err, "could not get ops")
	}

	count := 0

	for rows.Next() {
		var opID string
		var done bool
		if err := rows.Scan(&opID, &done); err != nil {
			return util.Annotate(err, "could not scan")
		}
		fmt.Printf("op: %s, done: %v\n", opID, done)
		count++
	}

	fmt.Printf("Count: %d\n", count)

	return nil
}

func (o *PQOperations) ReadOperation(ctx context.Context, id string) (*spb.Operation, error) {
	return nil, nil
}

func (o *PQOperations) NewOperation(ctx context.Context) (*spb.Operation, error) {
	return nil, nil
}
