package util

import (
	"fmt"
	"math/rand/v2"
)

func RandDigit(size int) string {
	var max uint64 = 1

	for range size {
		max *= 10
	}

	return fmt.Sprintf("%0*d", size, rand.Uint64()%max)
}
