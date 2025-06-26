package executors

import (
	"context"
	"sync"

	"github.com/oxia-db/nignly-testing/internal/assertions"
	"github.com/oxia-db/nignly-testing/internal/operations"
	"github.com/oxia-db/oxia/oxia"
	"github.com/pkg/errors"
)

var _ Executor = &oxiaExecutor{}

type oxiaExecutor struct {
	sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	syncClient oxia.SyncClient
}

func (o *oxiaExecutor) Close() error {
	o.cancel()
	o.Wait()

	return o.syncClient.Close()
}

func (o *oxiaExecutor) Execute(op operations.Operation) error {
	switch op.Type() {
	case operations.TypeSequence:
		sequenceOp := op.(*operations.SequenceOperation)
		return o.executeSequence(sequenceOp)
	}
	return nil
}

func (o *oxiaExecutor) executeSequence(op *operations.SequenceOperation) error {
	key := op.Key
	value := op.Value
	deltas := op.Deltas
	insertedKey, _, err := o.syncClient.Put(o.ctx, key, value, oxia.SequenceKeysDeltas(deltas...), oxia.PartitionKey(key))
	if err != nil {
		return err
	}
	assertion := op.Assertion
	if assertion == nil {
		return nil
	}
	if assertion.Key != nil && insertedKey != *assertion.Key {
		return errors.Wrapf(assertions.ErrAssertionFailed, "assertion key failed. expect: %v , actual %v", *assertion.Key, insertedKey)
	}
	return nil
}

type OxiaExecutorOptions struct {
	Address   string `json:"address" yaml:"address"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

func NewExecutor(ctx context.Context, option OxiaExecutorOptions) (Executor, error) {
	client, err := oxia.NewSyncClient(option.Address)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	return &oxiaExecutor{
		ctx:        ctx,
		cancel:     cancelFunc,
		syncClient: client,
	}, nil
}
