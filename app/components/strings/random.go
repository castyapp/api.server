package strings

import (
	"math/rand"
	"time"
)

func Random(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("0123456789")
	b := make([]rune, length)
	for i := range b {

		if letter := letters[rand.Intn(len(letters))]; letter != 0 {
			b[i] = letters[rand.Intn(len(letters))]
		}

		continue
	}
	return string(b)
}
