package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type logWriter struct {
	filename string
	fd       io.WriteCloser
	mu       *sync.Mutex
}

func newLogWriter(filename string) (*logWriter, error) {
	w := &logWriter{
		filename: filename,
		mu:       &sync.Mutex{},
	}
	return w, w.openFile()
}

func (w *logWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fd.Write(b)
}

func (w *logWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closeFile()
}

func (w *logWriter) Rotate() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.closeFile(); err != nil {
		return err
	}

	if err := w.openFile(); err != nil {
		return err
	}

	return nil
}

func (w *logWriter) openFile() error {
	fdnew, err := loadFile(w.filename)

	if err != nil {
		w.fd = os.Stderr
		return err
	}

	w.fd = fdnew

	return nil
}

func (w *logWriter) closeFile() error {
	if w.filename == "" {
		return nil
	}

	if err := w.fd.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}

	w.fd = os.Stderr

	return nil
}

func loadFile(filename string) (*os.File, error) {
	if filename == "" {
		return os.Stdout, nil
	}

	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("path %q error: %w", filename, err)
	}

	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}
