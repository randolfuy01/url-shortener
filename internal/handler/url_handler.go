package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/randolfuy01/url-shortener/internal/database"
	"github.com/randolfuy01/url-shortener/internal/shortener"
)

// Create User Handler
func Create_URL_Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Handle_Create_URL(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid method")
	}
}

func Handle_Create_URL(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("UserID"))
	if err != nil || userID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No UserID provided")
		return
	}
	url := database.URL{
		UserID:      userID,
		OriginalURL: r.URL.Query().Get("URL"),
	}
	Encryption_Choice := r.URL.Query().Get("Encryption")

	if url.OriginalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No provided URL for encryption")
		return
	}

	if Encryption_Choice == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No encryption choice provided")
		return
	}

	// Validate url
	if !Valid_URL(url.OriginalURL) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid URL")
		return
	}
	var encrpErr error
	// Shorten url
	switch Encryption_Choice {
	case "MD5":
		url.ShortURL, encrpErr = shortener.Encode(url.OriginalURL, shortener.Encryption_MD5)
	case "SHA256":
		url.ShortURL, encrpErr = shortener.Encode(url.OriginalURL, shortener.Encryption_SHA256)
	case "Cipher":
		url.ShortURL, encrpErr = shortener.Encode(url.OriginalURL, shortener.Encryption_Vigenere_Cipher)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid encrpytion")
		return
	}
	if encrpErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(w, "Failed during encrpytion step")
		return
	}

	url.ShortURL = url.ShortURL[:8]
	// Insert
	ok, dbErr := database.Insert_URL(url, database.Pool)
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failure writing into database")
		return
	}
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unsuccessful writing into database")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Successful url shortening")
}

func Valid_URL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

