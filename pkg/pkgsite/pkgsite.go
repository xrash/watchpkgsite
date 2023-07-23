package pkgsite

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

type Runtime struct {
	args            []string
	ctx             context.Context
	cancelCtx       context.CancelFunc
	killedCtx       context.Context
	killedCancelCtx context.CancelFunc
	cmd             *exec.Cmd
	stdout          *bytes.Buffer
	stderr          *bytes.Buffer
	done            bool
	err             error
	mu              *sync.Mutex
}

func Run(
	lgr zerolog.Logger,
	workdir string,
	addr string,
) (*Runtime, error) {

	args := []string{
		"-http",
		addr,
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	killedCtx, killedCancelCtx := context.WithCancel(context.Background())

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	lgr.Info().Str("args", strings.Join(args, ",")).Msg("running pkgsite with args")
	cmd := exec.CommandContext(ctx, "pkgsite", args...)
	cmd.Dir = workdir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	runtime := &Runtime{
		args:            args,
		ctx:             ctx,
		cancelCtx:       cancelCtx,
		killedCtx:       killedCtx,
		killedCancelCtx: killedCancelCtx,
		cmd:             cmd,
		stdout:          stdout,
		stderr:          stderr,
		mu:              &sync.Mutex{},
	}

	go func() {
		defer func() {
			runtime.mu.Lock()
			defer runtime.mu.Unlock()
			runtime.done = true
			runtime.killedCancelCtx()
		}()

		if err := cmd.Run(); err != nil {
			runtime.mu.Lock()
			defer runtime.mu.Unlock()
			runtime.err = err
			lgr.Error().Str("err", err.Error()).Msg("error from pkgsite command")
		}
	}()

	return runtime, nil
}

func (r *Runtime) Kill() error {
	// This will kill the process.
	r.cancelCtx()
	return nil
}

func (r *Runtime) Killed() <-chan struct{} {
	return r.killedCtx.Done()
}
