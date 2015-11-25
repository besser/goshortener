package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/besser/goshortener/url"
)

//region TYPES

type Headers map[string]string

//endregion TYPES

//region CONST AND VARS

var (
	port    int
	urlBase string
)

//endregion CONST AND VARS

//region MAIN FUNCTIONS

func init() {
	port = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", port)
}

func main() {
	http.HandleFunc("/api/shorten", Shortener)
	http.HandleFunc("/r/", Redirector)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

//endregion MAIN FUNCTIONS

//region PUBLIC FUNCIONS

func Redirector(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if url := url.Find(id); url != nil {
		http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
	}
}

func Shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		answerWith(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
		return
	}

	url, new, err := url.GetUrl(extractUrl(r))

	if err != nil {
		answerWith(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if new {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	shortUrl := fmt.Sprintf("%s/r/%s", urlBase, url.Id)
	answerWith(w, status, Headers{"Location":shortUrl})
}

//endregion PUBLIC FUNCIONS

//region PRIVATE FUNCIONS

func answerWith(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(status)
}

func extractUrl(r *http.Request) string	{
	url := make([]byte, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

//endregion PRIVATE FUNCIONS