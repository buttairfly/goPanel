package routine

import (
	"context"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

// GracefulExit tries to meet targetGoroutines num but waits gracePeriod otherwise it panics
func GracefulExit(ctx context.Context, targetGoroutines int, gracePeriod time.Duration, logger *zap.Logger) {
	go func() {
		<-ctx.Done() // wait until program should exit
		now := time.Now()
		gracefulWaitTime := now.Add(gracePeriod)
		tickChannel := time.Tick(200 * time.Millisecond)
		logger.Info("graceful exit started", zap.Time("until", gracefulWaitTime))
		for gracefulWaitTime.After(now) {

			select {
			case <-tickChannel:
				currentGoroutines := runtime.NumGoroutine()
				if currentGoroutines <= targetGoroutines {
					logger.Info("graceful exit", zap.Int("numGoroutines", currentGoroutines), zap.Int("targetGoroutines", targetGoroutines))
					// force positive exit
					os.Exit(0)
				}
				logger.Info("graceful wait", zap.Int("numGoroutines", currentGoroutines), zap.Int("targetGoroutines", targetGoroutines))
			}
			now = time.Now()
		}
		debug.SetTraceback("all")
		logger.Panic("graceful panics")
	}()
}
