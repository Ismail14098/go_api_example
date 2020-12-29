package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var correctDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(string)
	if ok {
		expTimeParsed, err := time.Parse(time.RFC3339, date)
		if err == nil && time.Now().Before(expTimeParsed) {
			return true
		}
	}
	return false
}

var correctLogin validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().Interface().(string)
	return len(text) > 5 && len(text) < 31
}

var correctPassword validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().Interface().(string)
	return len(text) > 7 && len(text) < 37
}

var correctName validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().Interface().(string)
	return len(text) > 1 && len(text) < 41
}

var correctTitle validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().Interface().(string)
	return len(text) > 4
}

var correctText validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().Interface().(string)
	return len(text) > 0
}

var correctCategoryName validator.Func = func(fl validator.FieldLevel) bool {
	text := fl.Field().Interface().(string)
	return len(text) > 4
}