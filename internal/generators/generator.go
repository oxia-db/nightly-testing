package generators

import (
	"io"

	"github.com/oxia-db/nignly-testing/internal/operations"
)

type Generator interface {
	io.Closer

	Name() string

	Next() operations.Operation
}
