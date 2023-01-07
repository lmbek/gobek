package utils

import (
	"github.com/NineNineFive/go-local-web-gui/utils/random"
	"testing"
)

func TestGetInt(test *testing.T) {
	// Set the random seed to a fixed value to make the test deterministic
	random.SetRandomSeed(0)

	for i := 0; i < 100; i++ { // modify amount of tries
		// Call the function that returns the value to be tested
		value := random.GetInt(10000, 10100)

		// Check if the value is within the range of 10000 and 10100
		if value < 10000 || value > 10100 {
			test.Errorf("Value returned is outside the range of 10000 and 10100. Value: %d", value)
		}
		//println(value) // uncomment to check values
	}

}
