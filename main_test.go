package main

import (
	"math"
	"testing"
)

func TestGetSwap(t *testing.T) {
	swap := getSwap()
	if math.IsNaN(float64(swap)) {
		t.Error("swap is NaN")
	}
}

