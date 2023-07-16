package doceng

import (
	"context"
	"fmt"
	"sync"

	"github.com/xrash/watchpkgsite/pkg/pkgsite"

	"github.com/rs/zerolog"
)

type Command struct {
	Kind string
}

type DocEngine struct {
	logger     zerolog.Logger
	workdir    string
	addr       string
	fatalErrCh chan<- error
}

func NewDocEngine(
	logger zerolog.Logger,
	workdir string,
	addr string,
	fatalErrCh chan<- error,
) *DocEngine {
	return &DocEngine{
		logger:     logger,
		workdir:    workdir,
		addr:       addr,
		fatalErrCh: fatalErrCh,
	}
}

func (e *DocEngine) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	commandsCh <-chan *Command,
) error {
	go e.background(ctx, wg, commandsCh)
	return nil
}

func (e *DocEngine) Stop() error {
	return nil
}

func (e *DocEngine) background(
	ctx context.Context,
	wg *sync.WaitGroup,
	commandsCh <-chan *Command,
) {
	defer wg.Done()
	defer e.logger.Info().Msg("background done")

	runtime, err := pkgsite.Run(e.logger, e.workdir, e.addr)
	if err != nil {
		e.fatalErrCh <- fmt.Errorf("error running pkgsite: %w", err)
		return
	}

	for {
		select {

		case <-ctx.Done():
			e.logger.Info().Msg("got context done")
			if err := runtime.Kill(); err != nil {
				e.logger.Error().
					Str("err", err.Error()).
					Msg("error killing runtime")
			}
			return

		case command := <-commandsCh:
			switch command.Kind {
			case "reload":
				e.logger.Info().
					Msg("got reload command")

				if err := runtime.Kill(); err != nil {
					e.fatalErrCh <- fmt.Errorf("error killing runtime: %w", err)
					return
				}

				runtime, err = pkgsite.Run(e.logger, e.workdir, e.addr)
				if err != nil {
					e.fatalErrCh <- fmt.Errorf("error running pkgsite: %w", err)
					return
				}

			}
		}
	}
}
