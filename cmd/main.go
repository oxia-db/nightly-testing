package main

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/oxia-db/nignly-testing/cmd/common"
	"github.com/oxia-db/nignly-testing/cmd/correctness"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	rootCmd = &cobra.Command{
		Use:   "night",
		Short: "Oxia nightly testing commands ",
		Long:  "Oxia nightly testing commands ",
	}
)

func init() {
	rootCmd.AddCommand(correctness.Cmd)
	rootCmd.PersistentFlags().StringVarP(&common.ExecutorOptions.Address, "address", "a", "127.0.0.1:6648", "oxia server address")
	rootCmd.PersistentFlags().StringVarP(&common.ExecutorOptions.Namespace, "namespace", "n", "default", "oxia executor namespace")
}

func main() {
	pprof.Do(context.Background(),
		pprof.Labels("component", "main"),
		func(ctx context.Context) {
			if _, err := maxprocs.Set(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if err := rootCmd.Execute(); err != nil {
				os.Exit(1)
			}
		})
}
