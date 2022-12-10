package something

import (
	"regexp"
	"strings"
)

func IsValidName(name string) bool {
	if name == "" {
		// return that it is not valid, as name cannot be empty
		return false
	}

	// only allow a-ø, A-Ø, and 0-9 (numbers)
	hasSymbols := regexp.MustCompile(`[^A-Øa-ø0-9 @_.,()]+`).MatchString(name)
	hasBackslash := strings.Contains(name, "\\")

	if hasSymbols || hasBackslash {
		// it has symbols so return false as it is not valid
		return false
	} else {
		// It seems valid so return true
		return true
	}
}

func GetValidityGuideLines() string {
	return "only symbols: a-z, A-Z, æøåÆØÅ, 0-9, @_.,() are allowed"
}
