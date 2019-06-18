package arduinocom

import "time"

type statsType string

const (
	ardoinoErrorType statsType = "Error arduino"
	printType        statsType = "Info         "
	latchType        statsType = "Latch        "
)

type stats struct {
	event     statsType
	timeStamp time.Time
	message   string
}
