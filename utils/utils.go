package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandString generates a n-char short code for the archived link
// This is a pretty simple implementation and I (@batuhan) am sure
// there are a handful of better ways to handle this.
// Source: https://stackoverflow.com/a/22892986
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type NullableTime struct {
	time.Time
}
func (t NullableTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	} else {
		return t.Time.MarshalJSON()
	}
}
