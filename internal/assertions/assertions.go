package assertions

import "errors"

var (
	ErrAssertionFailed = errors.New("assertion failed")
)

type Assertion struct {
	Key               *string
	Value             []byte
	VersionID         *int64
	CreateTimestamp   *uint64
	ModifiedTimestamp *uint64
	ModifiedCount     *uint64
	Ephemeral         *bool
}
