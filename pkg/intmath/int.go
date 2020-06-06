package intmath

import (
	"go.uber.org/zap"
)

// Rescale maps a value x on a in scale to out scale
func Rescale(x, inMin, inMax, outMin, outMax int) int {
	if (inMax - inMin) == 0 {
		zap.L().Error("input range is zero, return 0")
		return 0
	}
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// Constrain contrains the input x within its boundaries min and max
func Constrain(x, min, max int) int {
	if x > max {
		return max
	}
	if x < min {
		return min
	}
	return x
}

// Abs returns the absolute value of x.
// It ignores the case Abs(math.Minimum)
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Max returns the integer maximum from a and b
func Max(a, b int) int {
	if b < a {
		return a
	}
	return b
}

// Min returns the integer minimum from a and b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
