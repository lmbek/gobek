package net_test

import (
	"github.com/NineNineFive/go-local-web-gui/utils/net"
	"strconv"
	"testing"
)

func BenchmarkGetInt(benchmark *testing.B) {
	// Reset the timer before running the code
	benchmark.ResetTimer()

	// Run the code b.N times
	for i := 0; i < benchmark.N; i++ {
		for j := 11450; j < 11453; j++ {
			net.IsPortUsed("localhost", strconv.Itoa(j))
		}
	}
}

func TestPortIsUsed(test *testing.T) {
	for i := 11450; i < 11453; i++ {
		net.IsPortUsed("localhost", strconv.Itoa(i))
	}
}
