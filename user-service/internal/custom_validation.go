package internal

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		hasLowercase := regexp.MustCompile("[a-z]").MatchString(password)
		hasUppercase := regexp.MustCompile("[A-Z]").MatchString(password)
		hasNumber := regexp.MustCompile("[0-9]").MatchString(password)
		hasSpecial := regexp.MustCompile("[!@#$%^&*()_+\\-=\\[\\]{};':\"\\|,.<>/?`~]").MatchString(password)

		return hasLowercase && hasUppercase && hasNumber && hasSpecial
	})
}