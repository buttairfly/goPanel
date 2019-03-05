package device

import "time"

type statsType int

const (
	errorType statsType = iota
	printType
	latchType
)

type stats struct {
	event     statsType
	timeStamp time.Time
	message   string
}
