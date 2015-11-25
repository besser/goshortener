package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/besser/goshortener/url"
)

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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port)))
}