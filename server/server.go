package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/besser/goshortener/url"
)

//region TYPES

type Headers map[string]string

type Redirector struct {
	stats chan string
}

//endregion

//region CONST AND VARS

var (
	logOn   *bool
	port    *int
	urlBase string
)

//endregion

//region MAIN FUNCTIONS

func init() {
	port = flag.Int("p", 8888, "port")
	logOn = flag.Bool("l", true, "log on/off")

	flag.Parse()

	urlBase = fmt.Sprintf("http://localhost:%d", *port)
}

func Run() {
	stats := make(chan string)
	defer close(stats)
	go registerStatistics(stats)

	url.ConfigRepository(url.NewRepoMem())

	http.HandleFunc("/api/shorten", Shortener)
	http.HandleFunc("/api/stats/", Statistics)
	http.Handle("/r/", &Redirector{stats})

	toLog("Starting server listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

//endregion

//region REDIRECTOR PUBLIC METHODS

func (red *Redirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	getUrlAndExecute(w, r, func(url *url.Url) {
		http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
		red.stats <- url.Id // Recordind statistics
	})
}

//endregion

//region HTTP HANDLERS

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

	toLog("URL %s was shortened successfully to %s.", url.Destination, shortUrl)
}

func Statistics(w http.ResponseWriter, r *http.Request) {
	getUrlAndExecute(w, r, func(url *url.Url) {
		json, err := json.Marshal(url.Stats())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, string(json))
	})
}

//endregion

//region PRIVATE FUNCIONS

func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

func getUrlAndExecute(
	w http.ResponseWriter,
	r *http.Request,
	exec func(*url.Url),
) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if u := url.Find(id); u != nil {
		exec(u)
	} else {
		http.NotFound(w, r)
	}
}

func toLog(format string, values ...interface{}) {
	if *logOn {
		log.Printf(fmt.Sprintf("%s\n", format), values...)
	}
}

func registerStatistics(ids <-chan string) {
	for id := range ids {
		url.RegisterClick(id)
		toLog("Click successfully registered for %s.\n", id)
	}
}

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

//endregion
