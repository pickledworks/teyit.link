package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bund *Bundle

func InitBundle() {
	Bund = &i18n.Bundle{DefaultLanguage: language.Turkish}
	Bund.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Bund.MustLoadMessageFile("i18n/tr.toml")
	Bund.MustLoadMessageFile("i18n/en.toml")
}

func GetBundle() *Bundle {
	return Bund
}
