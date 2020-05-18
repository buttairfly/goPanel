package arduinocom

import "time"

type statType string

const (
	ArdoinoErrorStatType statType = "Error"
	PrintStatType        statType = "Info "
	LatchStatType        statType = "Latch"
)

// Stat marks an event which will get printed
type Stat struct {
	Event     statType
	TimeStamp time.Time
	Message   string
}
