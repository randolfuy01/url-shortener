package test

/*
Testing for validation package:
Includes:
	- Regex
TO DO:
	- Payload validation
*/

import (
	"testing"

	"github.com/randolfuy01/url-shortener/pkg"
)

func Test_Regex_Validation_1(t *testing.T) {
	test_val := "https://google.com/search=/"
	regex_exp := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`

	valid, err := pkg.Format_validation(test_val, regex_exp)

	if err != nil {
		t.Errorf("Regex Validation Test 1: Expected successful regex comparison, but failed for input: %s Error: %v", test_val, err)
	}

	if !valid {
		t.Errorf("Regex Validation Test 1: Expected successful regex comparison, received: False")
	}
}

// Case 2: Unsuccessful, nonmatching regex format
func Test_Regex_Validation_2(t *testing.T) {
	test_val := "https://google.com"
	regex_exp := `[` // Invalid regex

	valid, err := pkg.Format_validation(test_val, regex_exp)

	if err == nil {
		t.Errorf("Regex Validation Test 2: Expected regex error, but got none for regex: %s", regex_exp)
	}

	if valid {
		t.Errorf("Regex Validation Test 2: Expected failed match due to regex error, got true")
	}
}

// Case 3: Unsuccessful, nonmatching string format
func Test_Regex_Validation_3(t *testing.T) {
	test_val := "not a url"
	regex_exp := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`

	valid, err := pkg.Format_validation(test_val, regex_exp)

	if err != nil {
		t.Errorf("Regex Validation Test 3: Unexpected error: %v", err)
	}

	if valid {
		t.Errorf("Regex Validation Test 3: Expected validation to fail for invalid URL string")
	}
}

// Case 4: Unsuccessful, no val provided
func Test_Regex_Validation_4(t *testing.T) {
	test_val := ""
	regex_exp := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`

	valid, err := pkg.Format_validation(test_val, regex_exp)

	if err == nil {
		t.Errorf("Regex Validation Test 4: Expected error for empty input string")
	}

	if valid {
		t.Errorf("Regex Validation Test 4: Expected validation to fail for empty string")
	}
}

// Case 5: Unsuccessful, no regex provided
func Test_Regex_Validation_5(t *testing.T) {
	test_val := "https://google.com"
	regex_exp := ""

	valid, err := pkg.Format_validation(test_val, regex_exp)

	if err == nil {
		t.Errorf("Regex Validation Test 5: Expected error for empty regex string")
	}

	if valid {
		t.Errorf("Regex Validation Test 5: Expected validation to fail due to missing regex")
	}
}
