package handler

import (
	"fmt"
	"net/http"

	"github.com/randolfuy01/url-shortener/internal/database"
	"github.com/randolfuy01/url-shortener/internal/shortener"
)

// Create User Handler
func Create_User_Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Handle_Create_User(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid method")
	}
}

// TO DO: Refactor
func Handle_Create_User(w http.ResponseWriter, r *http.Request) {
	user := database.User{
		Name:     r.URL.Query().Get("name"),
		Password: r.URL.Query().Get("password"),
	}

	if user.Name == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, "No name provided")
		return
	}
	if user.Password == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, "No password provided")
		return
	}

	// Hash password
	hashed, err := shortener.Encode(user.Password, shortener.Encryption_SHA256)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(w, "Error hashing password")
		return
	}
	user.Password = hashed

	res, err := database.Execute_Query(user, database.Insert_User)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error creating user")
		return
	}
	if !res {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "User already exists")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Create new user")
}
