package validators

import (
	"net/url"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	RegisterValidator(&Validator{
		Name:                     "media_source_uri",
		CallValidationEvenIfNull: true,
		HandleFunc: func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			uri, err := url.ParseRequestURI(value)
			if err != nil {
				return false
			}
			if strings.Contains(uri.Host, "spotify") {
				return true
			}
			return true
		},
		TranslationRegister: func(ut ut.Translator) error {
			return ut.Add("media_source_uri", "MediaSourceUri is not valid", true)
		},
		TranslationFunc: func(ut ut.Translator, fe validator.FieldError) (t string) {
			t, _ = ut.T("media_source_uri", "MediaSourceUri")
			return
		},
	})
}
