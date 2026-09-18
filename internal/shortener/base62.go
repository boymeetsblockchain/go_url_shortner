// Package shortener encodes IDs into short Base62 codes.
package shortener

import "strings"

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func ToBase62(n int) string {
	if n == 0 {
		return string(base62Chars[0])
	}
	var sb strings.Builder
	for n > 0 {
		remainder := n % 62
		sb.WriteByte(base62Chars[remainder])
		n /= 62
	}
	// digits came out reversed, so flip them
	runes := []rune(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
