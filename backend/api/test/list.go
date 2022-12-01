package test

import (
	"errors"
)

// List Do something here like reading data from a file or database and return results
func List(isAllowed bool) ([]string, error) {
	list := []string{"this", "is", "some", "data"}
	if isAllowed {
		return list, nil
	} else {
		return []string{}, errors.New("not allowed to give list")
	}
}
