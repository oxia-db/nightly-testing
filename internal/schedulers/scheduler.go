package schedulers

import (
	"context"
	"io"
	"log/slog"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/oxia-db/nignly-testing/internal/executors"
	"github.com/oxia-db/nignly-testing/internal/generators"
	"go.uber.org/multierr"
)

type Scheduler interface {
	io.Closer
}

var _ Scheduler = &scheduler{}

type scheduler struct {
	*slog.Logger
	sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	generator generators.Generator
	executor  executors.Executor
}

func (s *scheduler) Close() error {
	s.Info("closing scheduler")
	defer s.Info("closed scheduler")
	s.cancel()
	s.Wait()

	err := s.generator.Close()
	return multierr.Append(err, s.executor.Close())
}

func (s *scheduler) schedule(ctx context.Context) {
	defer s.Done()

	s.Info("start scheduling the operations")
	defer s.Info("finished scheduling the operations")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			nextOperation := s.generator.Next()
			_ = backoff.RetryNotify(func() error {
				return s.executor.Execute(nextOperation)
			}, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
				if err != nil {
					s.Warn("unexpected exception when execute the operation.",
						slog.Any("error", err), slog.Duration("retry-after", duration))
					return
				}
			})
		}
	}
}

type SchedulerOptions struct {
	Generator generators.Generator
	Executor  executors.Executor
}

func NewScheduler(ctx context.Context, options SchedulerOptions) Scheduler {
	ctx, cancel := context.WithCancel(ctx)
	sche := &scheduler{
		Logger:    slog.With(slog.String("component", "scheduler"), slog.String("generator", options.Generator.Name())),
		ctx:       ctx,
		cancel:    cancel,
		generator: options.Generator,
		executor:  options.Executor,
	}

	sche.Add(1)
	go pprof.Do(ctx, pprof.Labels("component", "scheduler"), sche.schedule)

	return sche
}
