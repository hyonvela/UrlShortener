// Функция преобразования

package urlshortener

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	length   = 10
)

func Shorten(id uint32) string {
	var result [length]byte
	for i := 0; i < length; i++ {
		result[i] = alphabet[0]
	}

	var i int
	for id > 0 && i < length {
		rem := id % 63
		result[i] = alphabet[rem]
		id = id / 63
		i++
	}

	for j := 0; j < i/2; j++ {
		result[j], result[i-j-1] = result[i-j-1], result[j]
	}

	return string(result[:])
}
