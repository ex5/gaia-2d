package util

import (
	"math/rand"
)

func Roll(desiredStdDev float64, desiredMean float64) float64 {
	return rand.NormFloat64()*desiredStdDev + desiredMean
}

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
