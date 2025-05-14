package pkg

import (
	"errors"
	"regexp"
)

// regex validation for formatting
func Format_validation(val string, expression string) (bool, error) {

	// input validation
	if len(val) == 0 {
		return false, errors.New("no string provided")
	}

	// expression cannot be empty
	if expression == "" {
		return false, errors.New("regex pattern cannot be empty")
	}

	// Compile the regex to check validity
	compiled, err := regexp.Compile(expression)
	if err != nil {
		return false, err
	}

	// compiled regex for matching the string
	match := compiled.MatchString(val)

	// wrong format
	if !match {
		return false, nil
	}

	return true, nil
}
