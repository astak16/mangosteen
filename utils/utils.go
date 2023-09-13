package utils

import "crypto/rand"

func RandNumber(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	digits := make([]byte, len)
	for i := range b {
		digits[i] = b[i]%10 + 48
	}
	return string(digits), nil
}
