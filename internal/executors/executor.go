package executors

import (
	"io"

	"github.com/oxia-db/nignly-testing/internal/operations"
)

type Executor interface {
	io.Closer

	Execute(operation operations.Operation) error
}
