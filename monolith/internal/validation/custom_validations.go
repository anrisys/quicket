package validation

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)


func RegisterCustomValidation(v *validator.Validate) {
	_ = v.RegisterValidation("gttoday", func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if !ok {
			return false
		}
		nextDay := time.Now().AddDate(0, 0, 1)
		nextDayMidnight := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		return date.After(nextDayMidnight)
	})

	_ = v.RegisterValidation("payStatus", func(fl validator.FieldLevel) bool {
		status, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		normalized := strings.ToLower(strings.TrimSpace(status))
		return normalized == "success" || normalized == "failed"
	})

	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		/*
		// pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?])[A-Za-z\d!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]{8,}$`
		pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z\d]).{8,}$`
	
		matched, err := regexp.MatchString(pattern, password)
		if err != nil {
			return false
		}
		return matched
		*/
		hasLowercase := regexp.MustCompile("[a-z]").MatchString(password)
		hasUppercase := regexp.MustCompile("[A-Z]").MatchString(password)
		hasNumber := regexp.MustCompile("[0-9]").MatchString(password)
		hasSpecial := regexp.MustCompile("[!@#$%^&*()_+\\-=\\[\\]{};':\"\\|,.<>/?`~]").MatchString(password)

		return hasLowercase && hasUppercase && hasNumber && hasSpecial
	})
}