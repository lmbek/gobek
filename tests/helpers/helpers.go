package helpers

import (
	"fmt"
	"testing"
)

func StandardTestChecking(test *testing.T, result any, expected any) {
	if ResultIsExpected(result, expected) {
		PrintExpected(expected)
	} else {
		PrintNotExpected(test, result, expected)
	}
}

// ResultIsExpected - checks with an if statement if result is the same as expected
func ResultIsExpected(result any, expected any) bool {
	if result == expected {
		return true
	} else {
		return false
	}
}

func PrintExpected(expected any) {
	fmt.Printf("\t\tGot expected: %v \n", expected)
}

// PrintNotExpected check if result is the same as expected
func PrintNotExpected(test *testing.T, result any, expected any) {
	test.Errorf("\t\tExpected %v, but got %v", expected, result)
}

func PrintGotFatalError(test *testing.T, result any) {
	test.Fatalf("\t\tTest found fatal error: %v \n", result)
}
