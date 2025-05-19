package shortener

import "errors"

func Encode(url string, encoder func(string) (string, bool)) (string, error) {
	// Pass in any encryption algorithm which returns the encrypted string and boolean whether it was successful or not
	encrypted, ok := encoder(url)
	if !ok {
		return url, errors.New("unable to encode")
	}

	return encrypted, nil
}
