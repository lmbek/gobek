package utils_test

import (
	"github.com/NineNineFive/go-local-web-gui/utils"
	"reflect"
	"testing"
)

func BenchmarkToString(benchmark *testing.B) {
	// Reset the timer before running the code
	benchmark.ResetTimer()

	// Run the code b.N times
	for i := 0; i < benchmark.N; i++ {
		utils.IntegerToString(50)
	}
}

func TestIntegerToString(test *testing.T) {
	newString := utils.IntegerToString(50)
	if reflect.TypeOf(newString).Kind() != reflect.String {
		test.Error("not a string")
	}
	println(newString)
}
