package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/besser/goshortener/url"
)

//region TYPES

type Headers map[string]string

//endregion

//region CONST AND VARS

var (
	port    int
	urlBase string
	stats   chan string
)

//endregion

//region MAIN FUNCTIONS

func init() {
	port = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", port)
}

func main() {
	stats = make(chan string)
	defer close(stats)
	go registerStatistics(stats)

	url.ConfigRepository(url.NewRepoMem())

	http.HandleFunc("/api/shorten", Shortener)
	http.HandleFunc("/api/stats/", Statistics)
	http.HandleFunc("/r/", Redirector)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

//endregion

//region PUBLIC FUNCIONS

func Redirector(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if u := url.Find(id); u != nil {
		http.Redirect(w, r, u.Destination, http.StatusMovedPermanently)

		// Recordind statistics
		stats <- id
	} else {
		http.NotFound(w, r)
	}
}

func Shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWith(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
		return
	}

	url, new, err := url.GetUrl(extractUrl(r))

	if err != nil {
		respondWith(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if new {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	shortUrl := fmt.Sprintf("%s/r/%s", urlBase, url.Id)
	respondWith(w, status, Headers{
		"Location": shortUrl,
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", urlBase, url.Id),
	})
}

func Statistics(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if u := url.Find(id); u != nil {
		json, err := json.Marshal(u.Stats())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, string(json))
	} else {
		http.NotFound(w, r)
	}
}

//endregion

//region PRIVATE FUNCIONS

func respondWith(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(status)
}

func respondWithJSON(w http.ResponseWriter, reply string) {
	respondWith(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, reply)
}

func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

func registerStatistics(ids <-chan string) {
	for id := range ids {
		url.RegisterClick(id)
		fmt.Printf("Click successfully registered for %s.\n", id)
	}
}

//endregion
