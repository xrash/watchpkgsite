package giteng

import (
	"context"
	"fmt"
	"time"
	"sync"

	"github.com/xrash/watchpkgsite/pkg/git"

	"github.com/rs/zerolog"
)

type Event struct {
	Kind string
}

type GitEngine struct {
	logger     zerolog.Logger
	workdir    string
	interval   time.Duration
	fatalErrCh chan<- error
}

func NewGitEngine(
	logger zerolog.Logger,
	workdir string,
	interval time.Duration,
	fatalErrCh chan<- error,
) *GitEngine {
	return &GitEngine{
		logger:     logger,
		workdir:    workdir,
		interval:   interval,
		fatalErrCh: fatalErrCh,
	}
}

func (e *GitEngine) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	eventsCh chan<- *Event,
) error {
	go e.background(ctx, wg, eventsCh)
	return nil
}

func (e *GitEngine) Stop() error {
	return nil
}

func (e *GitEngine) background(
	ctx context.Context,
	wg *sync.WaitGroup,
	eventsCh chan<- *Event,
) {
	defer wg.Done()
	defer e.logger.Info().Msg("background done")

	runopts := &git.RunOptions{}
	if e.workdir != "" {
		runopts.Dir = e.workdir
	}

	onGitStatus := func(r *git.GitStatusResult) {
		e.logger.Info().
			Str("localBranch", r.LocalBranch).
			Str("remoteBranch", r.RemoteBranch).
			Str("syncState", string(r.SyncState)).
			Msg("git status")

		switch r.SyncState {

		case git.UpToDate:
			return

		case git.Behind:
			_, err := git.GitMerge(runopts, []string{r.RemoteBranch})
			if err != nil {
				e.fatalErrCh <- err
				return
			}
			eventsCh <- &Event{
				Kind: "update",
			}

		case git.Ahead:
			e.fatalErrCh <- fmt.Errorf("local branch shouldn't be ahead")
			return

		}
	}

	for {
		select {
		case <-ctx.Done():
			e.logger.Info().Msg("got context done")
			return
		case <-time.After(e.interval):

			e.logger.Info().Msg("git fetch begin")

			_, err := git.GitFetch(runopts)
			if err != nil {
				e.logger.Error().
					Str("err", err.Error()).
					Msg("error from git fetch")
				continue
			}

			e.logger.Info().Msg("git fetch end")

			e.logger.Info().Msg("git status start")

			r, err := git.GitStatus(runopts)
			if err != nil {
				e.logger.Error().
					Str("err", err.Error()).
					Msg("error from git status")
				continue
			}

			e.logger.Info().Msg("git status end")

			onGitStatus(r)
		}
	}
}
