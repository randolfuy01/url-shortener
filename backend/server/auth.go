package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// issueJWT creates an HS256 JWT token with standard claims
func issueJWT(userID int64, username string, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	now := time.Now().Unix()
	payload := map[string]interface{}{
		"sub":  fmt.Sprintf("%d", userID),
		"name": username,
		"iat":  now,
		"exp":  time.Now().Add(ttl).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	b64 := func(b []byte) string {
		return base64.RawURLEncoding.EncodeToString(b)
	}

	unsigned := b64(headerJSON) + "." + b64(payloadJSON)

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	sig := mac.Sum(nil)
	token := unsigned + "." + b64(sig)

	return token, nil
}

func hashData(data string) [32]byte {
	sum := sha256.Sum256([]byte(data))
	return sum
}

func ensureDB() {
	if pool == nil {
		initConnection()
	}
}
