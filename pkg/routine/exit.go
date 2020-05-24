package routine

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// DetectExit detects a interrupt or sigterm signal and closes the returned channel
func DetectExit(ctx context.Context) context.Context {
	cancelCtx, cancel := context.WithCancel(ctx)
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer cancel()

		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			// special treatment for os.Interrupt
		case syscall.SIGTERM:
			// special treatment for syscall.SIGTERM
		}
	}()
	return cancelCtx
}
