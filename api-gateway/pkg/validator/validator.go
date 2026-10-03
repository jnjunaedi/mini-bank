package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s any) map[string]string {
	err := validate.Struct(s)

	if err == nil {
		return nil
	}

	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {

			fieldName := fieldErr.Field()
			tag := fieldErr.Tag()

			switch tag {
			case "required":
				errorsMap[fieldName] = fieldName + " wajib diisi"
			case "numeric":
				errorsMap[fieldName] = fieldName + " hanya boleh diisi dengan angka"
			case "email":
				errorsMap[fieldName] = fieldName + " format email tidak valid"
			case "min":
				errorsMap[fieldName] = fieldName + " tidak boleh kurang dari " + fieldErr.Param()
			case "len":
				errorsMap[fieldName] = fieldName + " harus memiliki panjang tepat " + fieldErr.Param() + " karakter"
			default:
				errorsMap[fieldName] = fieldName + " tidak valid (" + tag + ")"
			}
		}
	}

	return errorsMap
}
