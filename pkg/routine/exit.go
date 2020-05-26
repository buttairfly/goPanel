package routine

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DetectExit detects a interrupt or sigterm signal and closes the returned channel
func DetectExit(ctx context.Context) context.Context {
	cancelCtx, cancel := context.WithCancel(ctx)
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer cancel()
		<-signalChannel
		log.New(os.Stdout, "\t", 0).Printf("%v shutdown", time.Now().Format("2006-01-02T15:04:05.000Z0700"))
	}()
	return cancelCtx
}
