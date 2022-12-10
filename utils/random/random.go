package random

import (
	"math/rand"
)

func SetRandomSeed(seed int64) {
	rand.Seed(seed)
}

func GetInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
