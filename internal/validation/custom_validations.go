package validation

import (
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
}