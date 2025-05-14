package handler

import (
	"fmt"
	"net/http"
)

// Test handler
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		helloHandle(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Invalid method")
	}
}

// Test handle
func helloHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello World!")
}
