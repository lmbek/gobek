package main

import (
	"testing"
)

func TestLaunchApp(test *testing.T) {
	var err error

	// TODO: let launchApp return an error, if it could not open app correctly
	launchApp()
	if err != nil {
		test.Errorf("failed")
	}

}
