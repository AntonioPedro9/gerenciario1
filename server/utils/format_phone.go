package utils

import (
	"errors"
	"regexp"
)

func FormatPhone(phone string) (string, error) {
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	if len(digits) == 10 {
		digits = digits[:2] + "9" + digits[2:]
	}

	if len(digits) != 11 {
		return "", errors.New("invalid phone number lenght")
	}

	return digits, nil
}
