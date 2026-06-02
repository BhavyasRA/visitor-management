package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(data any) string {
	err := validate.Struct(data)

	if err == nil {
		return ""
	}

	var messages []string

	for _, e := range err.(validator.ValidationErrors) {
		messages = append(messages, e.Field()+" is invalid")
	}

	return strings.Join(messages, ", ")
}
