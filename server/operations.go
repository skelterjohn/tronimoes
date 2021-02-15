package server

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "github.com/skelterjohn/tronimoes/server/proto"
)

type InMemoryOperations struct {
	opsMu sync.Mutex
	ops   map[string]*spb.Operation
}

func (o *InMemoryOperations) WriteOperation(ctx context.Context, op *spb.Operation) error {
	o.opsMu.Lock()
	defer o.opsMu.Unlock()
	if o.ops == nil {
		o.ops = map[string]*spb.Operation{}
	}
	o.ops[op.GetOperationId()] = op
	return nil
}

func (o *InMemoryOperations) ReadOperation(ctx context.Context, id string) (*spb.Operation, error) {
	o.opsMu.Lock()
	defer o.opsMu.Unlock()
	if op, ok := o.ops[id]; ok {
		return op, nil
	}
	return nil, status.Errorf(codes.NotFound, "no such operation")
}

func (o *InMemoryOperations) NewOperation(ctx context.Context) (*spb.Operation, error) {
	op := &spb.Operation{
		OperationId: uuid.New().String(),
	}
	if err := o.WriteOperation(ctx, op); err != nil {
		return nil, err
	}
	return op, nil
}
