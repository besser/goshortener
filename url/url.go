package url

import (
	"math/rand"
	"net/url"
	"time"
)

//region TYPES

type Url struct {
	Id              string
	CreationDate    time.Time
	Destination     string
}

type Repository interface {
	IdExists(id string) bool
	FindById(id string) *Url
	FindByUrl(url string) *Url
	Save(url Url) error
}

//endregion

//region CONST AND VARS

const (
	size    = 5
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

var repo Repository

//endregion

//region MAIN FUNCTIONS

func init() {
	rand.Seed(time.Now().UnixNano())
}

//endregion

//region PUBLIC FUNCIONS

func ConfigRepository(r Repository) {
	repo = r
}

func Find(id string) *Url {
	return repo.FindById(id)
}

func GetUrl(destiny string) (u *Url, new bool, err error) {
	if u = repo.FindByUrl(destiny); u != nil {
		return u, false, nil
	}

	if _, err = url.ParseRequestURI(destiny); err != nil {
		return nil, false, err
	}

	url := Url{generateId(), time.Now(), destiny}
	repo.Save(url)

	return &url, true, nil
}

//endregion

//region PRIVATE FUNCIONS

func generateId() string {
	newId := func() string {
		id := make([]byte, size)

		for i := range id {
			id[i] = symbols[rand.Intn(len(symbols))]
		}

		return string(id)
	}

	for {
		if id := newId(); !repo.IdExists(id) {
			return id
		}
	}
}

//endregion
