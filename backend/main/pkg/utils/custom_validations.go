package utils

import (
	"github.com/go-playground/validator"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
)

var validate *validator.Validate

func SetValidator(v *validator.Validate) {
	validate = v
}

func GetValidator() *validator.Validate {
	return validate
}

func ValidatePurl(fl validator.FieldLevel) bool {
	_, err := purlutils.PurlFromString(fl.Field().String())

	return err == nil
}
