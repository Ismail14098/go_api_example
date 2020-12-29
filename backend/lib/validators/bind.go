package validators

import "github.com/go-playground/validator/v10"

func Validate(v *validator.Validate){
	v.RegisterValidation("correctDate", correctDate)
	v.RegisterValidation("correctLogin", correctLogin)
	v.RegisterValidation("correctPassword", correctPassword)
	v.RegisterValidation("correctName", correctName)
	v.RegisterValidation("correctTitle", correctTitle)
	v.RegisterValidation("correctText", correctText)
	v.RegisterValidation("correctCategoryName", correctCategoryName)
}
