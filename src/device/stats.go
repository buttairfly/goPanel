package device

import "time"

type statsType int

const (
	ardoinoErrorType statsType = iota
	printType
	latchType
)

type stats struct {
	event     statsType
	timeStamp time.Time
	message   string
}
