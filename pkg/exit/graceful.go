package exit

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

// GracefulExit tries to meet targetGoroutines num but waits gracePeriod otherwise it panics
func GracefulExit(ctx context.Context, targetGoroutines int, gracePeriod time.Duration, tickPeriod time.Duration, logger *zap.Logger) {
	go func() {
		<-ctx.Done() // wait until program should exit
		now := time.Now()
		gracefulWaitTime := now.Add(gracePeriod)
		tickChannel := time.Tick(tickPeriod)
		logger.Info("graceful exit started", zap.Time("until", gracefulWaitTime))
		for gracefulWaitTime.After(now) {

			select {
			case <-tickChannel:
				currentGoroutines := runtime.NumGoroutine()
				gc := GetGoroutine(logger)
				if currentGoroutines <= targetGoroutines {
					logger.Info("graceful exit", zap.String("goroutines", fmt.Sprintf("%#v", gc)), zap.Int("targetGoroutines", targetGoroutines))
					// force positive exit
					os.Exit(0)
				}
				shortGoroutine := ""
				for _, g := range gc {
					if shortGoroutine != "" {
						shortGoroutine += ", "
					}
					shortGoroutine += g.ShortGoroutine()
				}
				logger.Info("graceful wait", zap.String("goroutines", shortGoroutine), zap.Int("targetGoroutines", targetGoroutines))
			}
			now = time.Now()
		}
		debug.SetTraceback("all")
		logger.Panic("graceful panics")
	}()
}
