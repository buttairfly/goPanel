package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Detect detects a interrupt or sigterm signal and closes the returned channel
func Detect() <-chan bool {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool)

	go func() {
		defer close(done)
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			// special treatment for os.Interrupt
		case syscall.SIGTERM:
			// special treatment for syscall.SIGTERM
		}
	}()
	return done
}