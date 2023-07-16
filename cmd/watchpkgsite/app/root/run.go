package root

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/xrash/watchpkgsite/pkg/doceng"
	"github.com/xrash/watchpkgsite/pkg/giteng"
)

func run(opts cliopts, args []string) int {
	exit := make(chan error, 8)
	ctx, cancelCtx := context.WithCancel(context.Background())
	wgEngines := &sync.WaitGroup{}
	wgEngines.Add(2)
	defer wgEngines.Wait()

	lgr, err := makeLogger(opts.logger, opts.logfile)
	if err != nil {
		fmt.Printf("error making logger: %v", err)
		return 1
	}

	docengCmds := make(chan *doceng.Command, 1024)
	gitengEvents := make(chan *giteng.Event, 1024)

	go func() {
		for event := range gitengEvents {
			switch event.Kind {
			case "update":
				docengCmds <- &doceng.Command{
					Kind: "reload",
				}
			}
		}
	}()

	docEngine := doceng.NewDocEngine(
		lgr.With().Str("source", "doc_engine").Logger(),
		opts.workdir,
		opts.addr,
		exit,
	)

	if err := docEngine.Start(ctx, wgEngines, docengCmds); err != nil {
		fmt.Printf("error running doc engine: %v", err)
		return 1
	}

	defer func() {
		if err := docEngine.Stop(); err != nil {
			lgr.Error().
				Str("err", err.Error()).
				Msg("error stopping doc engine")
		}
	}()

	gitEngine := giteng.NewGitEngine(
		lgr.With().Str("source", "git_engine").Logger(),
		opts.workdir,
		opts.interval,
		exit,
	)

	if err := gitEngine.Start(ctx, wgEngines, gitengEvents); err != nil {
		fmt.Printf("error running git engine: %v", err)
		return 1
	}

	defer func() {
		if err := gitEngine.Stop(); err != nil {
			lgr.Error().
				Str("err", err.Error()).
				Msg("error stopping git engine")
		}
	}()

	sigint := make(chan os.Signal, 8)
	sighup := make(chan os.Signal, 8)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sighup, syscall.SIGHUP)

	for {
		select {

		case err := <-exit:
			if err == nil {
				lgr.Info().Msg("got exit signal with no error")
				cancelCtx()
				return 0
			}
			lgr.Error().Err(err).Msg("got fatal err")
				cancelCtx()
			return 1

		case <-sigint:
			lgr.Info().Msg("got sigint, exiting")
			cancelCtx()
			return 0

		case <-sighup:
			// @TODO reload
		}
	}

	return 0
}

func makeLogger(enabled bool, filename string) (*zerolog.Logger, error) {
	if !enabled {
		logger := zerolog.New(io.Discard)
		return &logger, nil
	}

	writer, err := makeWriter(filename)
	if err != nil {
		return nil, fmt.Errorf("error making writer: %v", err)
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger, nil
}

func makeWriter(filename string) (io.WriteCloser, error) {
	if filename == "" {
		return os.Stdout, nil
	}

	dirname := filepath.Dir(filename)
	if err := mkdir(dirname); err != nil {
		return nil, fmt.Errorf("error creating dir %v: %w", dirname, err)
	}

	writer, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	return writer, nil
}

func mkdir(dirname string) error {
	return os.MkdirAll(dirname, os.ModePerm)
}
