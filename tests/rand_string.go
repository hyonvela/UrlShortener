package tests

import (
	"math/rand/v2"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func RandString() string {
	len := int(rand.UintN(250)) + 1
	var b strings.Builder

	for i := 0; i < len; i++ {
		b.WriteByte(alphabet[rand.UintN(63)])
	}

	return b.String()
}
