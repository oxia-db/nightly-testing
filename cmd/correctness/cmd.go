package correctness

import (
	"github.com/oxia-db/nignly-testing/cmd/correctness/sequence"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "correctness",
		Short: "correctness verification",
		Long:  `correctness verification`,
	}
)

func init() {
	Cmd.AddCommand(sequence.Cmd)
}
