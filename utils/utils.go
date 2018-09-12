package utils

import (
	"errors"
	"math/rand"
	"net/url"
	"strings"
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

var (
	ErrorUrlEmptyUrl = errors.New("empty url")
	ErrorUrlInvalid = errors.New("invalid url")
	ErrorUrlSchemeInvalid = errors.New("invalid url scheme, must be http or https")
	ErrorUrlBlocked = errors.New("invalid url")
)

func ValidateUrl(urlStr string) error {
	if urlStr == "" {
		return ErrorUrlEmptyUrl
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return ErrorUrlInvalid
	} else if u.Scheme == "" || u.Host == "" {
		return ErrorUrlInvalid
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return ErrorUrlSchemeInvalid
	} else if strings.Contains(u.Host, "teyit.link") {
		return ErrorUrlBlocked
	}

	return nil
}
