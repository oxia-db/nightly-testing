package sequence

import (
	"io"
	"os"
	"syscall"

	"github.com/google/uuid"
	"github.com/opendss/toolkit/pkg/computing"
	"github.com/oxia-db/nignly-testing/cmd/common"
	"github.com/oxia-db/nignly-testing/internal/executors"
	"github.com/oxia-db/nignly-testing/internal/generators"
	"github.com/oxia-db/nignly-testing/internal/schedulers"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "sequence",
		Short: "put with sequence correctness verification",
		Long:  `put with sequence correctness verification`,
		RunE:  exec,
	}
	generatorOption = generators.SequenceGeneratorOptions{}
)

func init() {
	Cmd.Flags().StringVarP(&generatorOption.Key, "key", "k", uuid.NewString(), "the key of sequence")
	Cmd.Flags().IntVarP(&generatorOption.ValueSize, "message-size", "s", 1024, "the size of sequence in bytes")
}

func exec(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	generator := generators.NewSequenceGenerator(generatorOption)
	executor, err := executors.NewExecutor(ctx, common.ExecutorOptions)
	if err != nil {
		return err
	}
	computing.RunProcessWaitSignal(func() (io.Closer, error) {
		scheduler := schedulers.NewScheduler(ctx, schedulers.SchedulerOptions{
			Generator: generator,
			Executor:  executor,
		})
		return scheduler, nil
	}, os.Interrupt, syscall.SIGTERM)
	return nil
}
