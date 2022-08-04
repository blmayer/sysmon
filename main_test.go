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

func Test_getNet(t *testing.T) {
	netIn, netOut := getNet()
	if netIn < 0 || netOut < 0 {
		t.Error("negative numbers")
	}
}
