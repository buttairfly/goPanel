package exit

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// DetectSignal detects a interrupt or sigterm signal and closes the returned channel
func DetectSignal(ctx context.Context, logger *zap.Logger) context.Context {
	cancelCtx, cancel := context.WithCancel(ctx)
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer cancel()
		<-signalChannel
		logger.Info("\n")
		logger.Info("shutdown detedcted")
	}()
	return cancelCtx
}
