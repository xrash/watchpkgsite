package pkgsite

import (
	"bytes"
	"context"
	"os/exec"
	"sync"

	"github.com/rs/zerolog"
)

type Runtime struct {
	args      []string
	ctx       context.Context
	cancelCtx context.CancelFunc
	cmd       *exec.Cmd
	stdout    *bytes.Buffer
	stderr    *bytes.Buffer
	done      bool
	err       error
	mu        *sync.Mutex
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

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, "pkgsite", args...)
	cmd.Dir = workdir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	runtime := &Runtime{
		args:      args,
		ctx:       ctx,
		cancelCtx: cancelCtx,
		cmd:       cmd,
		stdout:    stdout,
		stderr:    stderr,
		mu:        &sync.Mutex{},
	}

	go func() {
		defer func() {
			runtime.mu.Lock()
			defer runtime.mu.Unlock()
			runtime.done = true
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
	r.cancelCtx()
	return nil
}
