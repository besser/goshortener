package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/besser/goshortener/url"
)

type Headers map [ string ] string

var (
  	port 	int
	urlBase	string
)

func init() {
	port = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", port)
}

func main() {
	http.HandleFunc("/api/short", Shortener)
	http.HandleFunc("/r/", Redirect)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func Shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		answerWith(w, http.StatusMethodNotAllowed, Headers { "Allow": "POST" })
		return
	}
}

func answerWith(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	
	w.WriteHeader(status)
}