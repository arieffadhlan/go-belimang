package utils

import (
    "net/url"
    "strings"

    "github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
  v.RegisterValidation("validUrl", func(fl validator.FieldLevel) bool {
    str := fl.Field().String()
    if str == "" {
      return false
    }

    u, err := url.ParseRequestURI(str)
    if err != nil {
      return false
    }

    return u.Scheme != "" && u.Host != "" && strings.Contains(u.Host, ".")
  })
}
