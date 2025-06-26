package generators

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/oxia-db/nignly-testing/internal/assertions"
	"github.com/oxia-db/nignly-testing/internal/operations"
)

var _ Generator = &sequenceGenerator{}

type sequenceGenerator struct {
	SequenceGeneratorOptions

	currentValueCounter atomic.Uint64
}

func (s *sequenceGenerator) Close() error {
	return nil
}

func (s *sequenceGenerator) Name() string {
	return "sequence"
}

func (s *sequenceGenerator) Next() operations.Operation {
	currentCounter := s.currentValueCounter.Add(1)

	payload := make([]byte, s.ValueSize)
	binary.BigEndian.PutUint64(payload, currentCounter)

	deltas := make([]uint64, 3)
	deltas[0] = 1
	deltas[1] = 2
	deltas[2] = 3

	assertKey := fmt.Sprintf("%v-%020d-%020d-%020d", s.Key, currentCounter, currentCounter*2, currentCounter*3)
	return &operations.SequenceOperation{
		Key:    s.Key,
		Value:  payload,
		Deltas: deltas,
		Assertion: &assertions.Assertion{
			Key: &assertKey,
		},
	}
}

type SequenceGeneratorOptions struct {
	Key       string
	ValueSize int
}

func NewSequenceGenerator(options SequenceGeneratorOptions) Generator {
	return &sequenceGenerator{
		SequenceGeneratorOptions: options,
		currentValueCounter:      atomic.Uint64{},
	}
}
