package validators

import (
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	RegisterValidator(&Validator{
		Name:                     "access",
		CallValidationEvenIfNull: true,
		HandleFunc: func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if value == "" {
				return true
			}
			switch value {
			case "1", "2", "3":
				return true
			}
			return false
		},
		TranslationRegister: func(ut ut.Translator) error {
			return ut.Add("access", "The access field format is invalid", true)
		},
		TranslationFunc: func(ut ut.Translator, fe validator.FieldError) (t string) {
			t, _ = ut.T("access", "access")
			return
		},
	})
}
