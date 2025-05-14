package shortener

import "errors"

// TO DO: Implement the encoding function and create tests for the function
// Using different encoding methods:
//   - MD5
//   - SHA-256
//   - xxHash
//   - Vigerence Cipher
func Encode(url string, encoder func(string) (string, bool)) (string, error) {
	// Pass in any encryption algorithm which returns the encrypted string and boolean whether it was successful or not
	encrypted, ok := encoder(url)
	if !ok {
		return url, errors.New("unable to generated shortened function")
	}

	return encrypted[:8], nil
}
