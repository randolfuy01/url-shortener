package handler

import (
	"fmt"
	"net/http"

	"github.com/randolfuy01/url-shortener/internal/database"
)

// Create User Handler
func Create_User_Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Create_Handle(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid method")
	}
}

// TO DO: Refactor
func Create_Handle(w http.ResponseWriter, r *http.Request) {
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
