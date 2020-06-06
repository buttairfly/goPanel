package palette

import "fmt"

// NoKeyColorFoundError returns a error at pos
func NoKeyColorFoundError(pos float64) error {
	return fmt.Errorf("no keycolor for palette found at %v", pos)
}
