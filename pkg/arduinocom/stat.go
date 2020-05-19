package arduinocom

import "time"

type statType string

const (
	// ArdoinoErrorStatType is the error category from the arduino
	ArdoinoErrorStatType statType = "Error"
	// PrintStatType is the into category from the arduino
	PrintStatType statType = "Info "
	// LatchStatType is the current status category from the arduino
	LatchStatType statType = "Latch"
)

// Stat marks an event which will get printed
type Stat struct {
	Event     statType
	TimeStamp time.Time
	Message   string
}
