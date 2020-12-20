package routine

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// DetectExit detects a interrupt or sigterm signal and closes the returned channel
func DetectExit(ctx context.Context, logger *zap.Logger) context.Context {
	cancelCtx, cancel := context.WithCancel(ctx)
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer cancel()
		<-signalChannel
		logger.Info("shutdown detedcted")
	}()
	return cancelCtx
}
