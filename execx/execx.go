package execx

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunInteractiveCapture(maxCaptureBytes int, name string, args ...string) (string, error) {
	if maxCaptureBytes <= 0 {
		maxCaptureBytes = 256 * 1024
	}
	var capBuf limitedBuffer
	capBuf.max = maxCaptureBytes

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = io.MultiWriter(os.Stdout, &capBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &capBuf)
	err := cmd.Run()
	return capBuf.String(), err
}

func RunCombined(timeout time.Duration, name string, args ...string) (string, error) {
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

type limitedBuffer struct {
	mu  sync.Mutex
	b   []byte
	max int
}

func (w *limitedBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(p) >= w.max {
		w.b = append(w.b[:0], p[len(p)-w.max:]...)
		return len(p), nil
	}
	if len(w.b)+len(p) > w.max {
		drop := len(w.b) + len(p) - w.max
		w.b = append(w.b[drop:], p...)
		return len(p), nil
	}
	w.b = append(w.b, p...)
	return len(p), nil
}

func (w *limitedBuffer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return string(w.b)
}
