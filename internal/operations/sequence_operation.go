package operations

import "github.com/oxia-db/nignly-testing/internal/assertions"

var _ Operation = &SequenceOperation{}

type SequenceOperation struct {
	Key    string
	Value  []byte
	Deltas []uint64

	Assertion *assertions.Assertion
}

func (s *SequenceOperation) Type() Type {
	return TypeSequence
}
