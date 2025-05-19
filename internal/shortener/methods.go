package shortener

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
)

// md5 encryption algorithm
func Encryption_MD5(url string) (string, bool) {

	// String is empty
	if len(url) == 0 {
		return url, false
	}
	// Initialize the new hash
	hash := md5.New()
	_, err := io.WriteString(hash, url)

	// Hashing went wrong
	if err != nil {
		return url, false
	}

	// Encode and send
	hashed := hash.Sum(nil)
	encoded := hex.EncodeToString(hashed)
	return encoded, true
}

// sha256 encrpytion algorithm
func Encryption_SHA256(url string) (string, bool) {
	// Initialize the new hash
	hash := sha256.New()
	_, err := io.WriteString(hash, url)

	// Hashing went wrong
	if err != nil {
		return url, false
	}

	// Encode and send
	hashed := hash.Sum(nil)
	encoded := hex.EncodeToString(hashed)
	return encoded, true
}

// vigenere cipher algorithm
func Encryption_Vigenere_Cipher(url string) (string, bool) {
	// Import key
	err := godotenv.Load(".env")

	// Failed loading environment
	if err != nil {
		log.Fatal("unable to load environment")
		return url, false
	}

	key := os.Getenv("VIGENERE_KEY")
	if len(key) == 0 {
		log.Fatal("no key provided for vigenere cipher")
		return url, false
	}

	// Normalize to uppercase
	url = strings.ToUpper(url)
	key = strings.ToUpper(key)

	// Cipher
	var encryption []rune
	url_runes := []rune(url)
	key_runes := []rune(key)

	j := 0
	for i := range url_runes {

		// Skip non-alphabetic characters
		if !unicode.IsLetter(url_runes[i]) {
			encryption = append(encryption, url_runes[i])
			continue
		} else if j < len(key_runes) {
			// Encrypt only alphabetic characters
			character := (url_runes[i]-'A'+key_runes[j]-'A')%26 + 'A'
			encryption = append(encryption, character)
			j++
		}

		if j == len(key_runes) {
			j = 0
		}
	}

	return string(encryption), true
}
