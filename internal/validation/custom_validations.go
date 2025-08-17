package validation

import (
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
}