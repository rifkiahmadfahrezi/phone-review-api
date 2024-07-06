package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

func IsValidUrl(str string) bool {
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	// jika tidak ada http atau https
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if !strings.Contains(u.Host, ".") {
		return false
	}

	return true
}

func CustomBindError(err error) string {
	var errorMessages []string

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("Field %s harus diisi", err.Field()))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("Field %s harus merupakan email yang valid", err.Field()))
			case "url":
				errorMessages = append(errorMessages, fmt.Sprintf("Field %s harus merupakan url yang valid", err.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("Field %s harus memiliki nilai minimal %s", err.Field(), err.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("Field %s harus memiliki nilai maksimal %s", err.Field(), err.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Field %s invalid", err.Field()))
			}
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	return strings.Join(errorMessages, ", ")
}
