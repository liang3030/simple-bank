package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/liang3030/simple-bank/util"
)

// declares a variable validCurrency of type validator.Func.
// validator.Func is a function signature defined by the validation package
var validCurrency validator.Func = func(field validator.FieldLevel) bool {
	/*
		field.Field() retrieves the actual field (value) being validated.
		.Interface() returns the field's value as an interface{}, which is a Go type that can represent any type.
		(string) is the type assertion. This attempts to convert the interface{} to a string.
	*/
	if currency, ok := field.Field().Interface().(string); ok {
		return util.IsValidCurrency(currency)
	}

	return false
}
