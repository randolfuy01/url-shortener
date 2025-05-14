package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/randolfuy01/url-shortener/internal/handler"
)

func main() {
	http.HandleFunc("/hello-word", handler.HelloHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
