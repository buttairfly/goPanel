package helper

import "log"

func IntMap(x, inMin, inMax, outMin, outMax int) int {
	if (inMax - inMin) == 0 {
		log.Printf("Error: input range is zero, return 0")
		return 0
	}
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func IntConstrain(x, min, max int) int {
	if x > max {
		return max
	}
	if x < min {
		return min
	}
	return x
}
