package utils

import "math"

func CustomRound(value float64) int {
	if math.Mod(value, 1) >= 0.5 {
		return int(math.Ceil(value))
	}
	return int(math.Floor(value))
}
