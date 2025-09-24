package test

/*
Testing for the encoders to ensure that hashes are being performed correctly across all available functions.
Should generate the correct answers based on the correct hashing functions
TO DO:
- Implement testing for: vigerance cipher
*/

import (
	"testing"

	"url-shortener.com/m/pkg"
)

// Case 1: Proper Encryption
func Test_Encryption_test_md5_1(t *testing.T) {
	input := "Hello_world"
	expected := "a2e1403220fc3fba0c291507a68e438f"
	output, ok := pkg.Encryption_MD5(input)

	if !ok {
		t.Errorf("MD5 Encryption Test 1: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`MD5 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 2: Large String
func Test_Encryption_test_md5_2(t *testing.T) {
	input := "This MD5 hash generator is useful for encoding passwords, credit cards numbers and other sensitive date into MySQL, Postgress or other databases. PHP programmers, ASP programmers and anyone developing on MySQL, SQL, Postgress or similar should find this online tool an especially handy resource."
	expected := "d1921aa0ca3c1146a01520c04e6caa9e"
	output, ok := pkg.Encryption_MD5(input)

	if !ok {
		t.Errorf("MD5 Encryption Test 2: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`MD5 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 3: Empty String
func Test_Encryption_test_md5_3(t *testing.T) {
	input := ""
	expected := ""
	output, ok := pkg.Encryption_MD5(input)

	if ok {
		t.Errorf("MD5 Encryption Test 3: Expected failed encryption, but encryption successful for input: %s", input)
	}

	if output != expected {
		t.Errorf(`MD5 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 4: Single Character
func Test_Encryption_test_md5_4(t *testing.T) {
	input := "F"
	expected := "800618943025315f869e4e1f09471012"
	output, ok := pkg.Encryption_MD5(input)

	if !ok {
		t.Errorf("MD5 Encryption Test 4: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`MD5 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 1: Proper Encryption of a string
func Test_Encryption_test_sha256_1(t *testing.T) {
	input := "Hello_world"
	expected := "af4e61dac3e037a684eedb5d9000f9cd3df89d911cc8e2354f73d3f6552b0357"
	output, ok := pkg.Encryption_SHA256(input)

	if !ok {
		t.Errorf("SHA-256 Encryption Test 1: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`SHA-256 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 2: Large String
func Test_Encryption_test_sha256_2(t *testing.T) {
	input := "This MD5 hash generator is useful for encoding passwords, credit cards numbers and other sensitive date into MySQL, Postgress or other databases. PHP programmers, ASP programmers and anyone developing on MySQL, SQL, Postgress or similar should find this online tool an especially handy resource."
	expected := "c12f0fa190b53aa95eb22006ca3f3a4bcb100af5b50d42170da93e946f1ebea4"
	output, ok := pkg.Encryption_SHA256(input)

	if !ok {
		t.Errorf("SHA-256 Encryption Test 2: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`SHA-256 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 3: Empty String
func Test_Encryption_test_sha256_3(t *testing.T) {
	input := ""
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	output, ok := pkg.Encryption_SHA256(input)

	if !ok {
		t.Errorf("SHA-256 Encryption Test 3: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`SHA-256 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Case 4: Single Character
func Test_Encryption_test_sha256_4(t *testing.T) {
	input := "F"
	expected := "f67ab10ad4e4c53121b6a5fe4da9c10ddee905b978d3788d2723d7bfacbe28a9"
	output, ok := pkg.Encryption_SHA256(input)

	if !ok {
		t.Errorf("SHA-256 Encryption Test 4: Expected successful encryption, but encryption failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`SHA-256 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Edge Case: Unicode input
func Test_Encryption_test_md5_unicode(t *testing.T) {
	input := "こんにちは"
	expected := "c0e89a293bd36c7a768e4e9d2c5475a8"
	output, ok := pkg.Encryption_MD5(input)

	if !ok {
		t.Errorf("MD5 Unicode Test: Expected successful encryption, but failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`MD5 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}

// Edge Case: Unicode input for SHA-256
func Test_Encryption_test_sha256_unicode(t *testing.T) {
	input := "こんにちは"
	expected := "125aeadf27b0459b8760c13a3d80912dfa8a81a68261906f60d87f4a0268646c"
	output, ok := pkg.Encryption_SHA256(input)

	if !ok {
		t.Errorf("SHA-256 Unicode Test: Expected successful encryption, but failed for input: %s", input)
	}

	if output != expected {
		t.Errorf(`SHA-256 Encryption of "%s" = %s, want %s`, input, output, expected)
	}
}
