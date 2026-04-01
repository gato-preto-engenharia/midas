/*
Package midas implement lifecycle utilities for go applications

Internally midas uses [slog] for logging, for more information see [SetupSlog] for
a way to use the default logging config

When to use each runner functions?
  - Use [Run] when you just need to start the [Runnable] and wait for it to finish
  - Use [Go] for long running processes with fine-grained control over cancellation, such as HTTP servers and background workers
  - Use [Supervise] for long running processes with managed cancellation and graceful shutdown
*/
package midas

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Runnable is midas primary interface, it represents a process that can have its lifecycle
// managed by midas functions
type Runnable interface {
	// Run is the [Runnable] midas entrypoint
	Run(context.Context, Config) error
}

// Run Executes the runnable
func Run(ctx context.Context, cfg Config, runnable Runnable) error {
	slog.Info("starting runnable")

	if err := runnable.Run(ctx, cfg); err != nil {
		slog.Error("sunnable finished with error", "error", err.Error())

		return err
	}

	slog.Info("runnable finished without error")

	return nil
}

// Go starts the runnable in a goroutine, returning a [context.CancelFunc] for manual shutdown,
// it is recommended to use [Supervise] if you do not need fine-grained control
// over the runnable shutdown
func Go(ctx context.Context, cfg Config, runnable Runnable) context.CancelFunc {
	cancelCtx, cancelFn := context.WithCancel(ctx)

	slog.Info("starting runnable in goroutine")

	go func() {
		if err := runnable.Run(cancelCtx, cfg); err != nil {
			slog.Error("runnable finished with error", "error", err.Error())
		} else {
			slog.Info("runnable finished without error")
		}
	}()

	return cancelFn
}

// Supervise Is midas runnable orchestrator, it spawns all the runnables with [Go], keep them alive
// and gracefully handle their shutdown and cancellation process. Available variables:
//   - MIDAS_CANCEL_WAIT_DURATION: control the time to wait after calling all runnable cancel functions, default to 100ms
func Supervise(ctx context.Context, cfg Config, runnables ...Runnable) {
	if ctx.Err() != nil {
		slog.Error("trying to supervise with faulty context, exiting supervise", "error", ctx.Err())
		return
	}

	cancels := make([]context.CancelFunc, len(runnables))
	for i, runnable := range runnables {
		cancels[i] = Go(ctx, cfg, runnable)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-sc // Hold app alive until a termination signal

	slog.Info("cancelling contexts")

	for _, cancel := range cancels {
		cancel()
	}

	time.Sleep(cfg.GetDuration("MIDAS_CANCEL_WAIT_DURATION", 100*time.Millisecond))
}

type wrapper struct {
	fn func(context.Context, Config) error
}

func (w *wrapper) Run(ctx context.Context, cfg Config) error {
	return w.fn(ctx, cfg)
}

// WrapFunc Wraps a function into a [Runnable]
func WrapFunc(fn func(context.Context, Config) error) Runnable {
	return &wrapper{
		fn: fn,
	}
}
