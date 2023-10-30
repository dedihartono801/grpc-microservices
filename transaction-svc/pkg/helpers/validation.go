package helpers

import (
	"regexp"
)

func ValidateInputSpecific(input string) bool {
	// Mencocokkan input dengan ekspresi reguler
	regex := regexp.MustCompile("^[a-zA-Z0-9.@]+$")
	return regex.MatchString(input)
}

func ValidateInputCommon(input string) bool {
	// Mencocokkan input dengan ekspresi reguler
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return regex.MatchString(input)
}
