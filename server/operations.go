package server

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "github.com/skelterjohn/tronimoes/server/proto"
)

type InMemoryOperations struct {
	opsMu sync.Mutex
	ops   map[string]*tpb.Operation
}

func (o *InMemoryOperations) WriteOperation(ctx context.Context, op *tpb.Operation) error {
	o.opsMu.Lock()
	defer o.opsMu.Unlock()
	if o.ops == nil {
		o.ops = map[string]*tpb.Operation{}
	}
	o.ops[op.OperationId] = op
	return nil
}

func (o *InMemoryOperations) ReadOperation(ctx context.Context, id string) (*tpb.Operation, error) {
	o.opsMu.Lock()
	defer o.opsMu.Unlock()
	if op, ok := o.ops[id]; ok {
		return op, nil
	}
	return nil, status.Errorf(codes.NotFound, "no such operation")
}

func (o *InMemoryOperations) NewOperation(ctx context.Context) (*tpb.Operation, error) {
	op := &tpb.Operation{
		OperationId: uuid.New().String(),
	}
	if err := o.WriteOperation(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}
