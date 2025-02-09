// Функция преобразования

package urlshortener

import (
	"strings"
)

const ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
const BASE = 63

func Shorten(id uint32) string {
	var (
		nums []uint32
		b    strings.Builder
		num  = id
	)

	for num > 0 {
		nums = append(nums, num%BASE)
		num /= BASE
	}

	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}

	for _, num := range nums {
		b.WriteString(string(ALPHABET[num]))
	}

	for b.Len() < 10 {
		b.WriteString(string(ALPHABET[0]))
	}

	if b.Len() > 10 {
		return b.String()[:10]
	}

	return b.String()
}
