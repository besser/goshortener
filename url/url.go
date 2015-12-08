package url

import (
	"math/rand"
	"net/url"
	"time"
)

//region TYPES

type Stats struct {
    Url     *Url    `json:"url"`
    Clicks  int     `json:"clicks"`
}

type Url struct {
	Id           string     `json:"id"`
	CreationDate time.Time  `json:"creationDate"`
	Destination  string     `json:"destination"`
}

//endregion

//region INTERFACES

type Repository interface {
    GetClicks(id string) int
	FindById(id string) *Url
	FindByUrl(url string) *Url
	IdExists(id string) bool
	RegisterClick(id string)
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

//region INIT

func init() {
	rand.Seed(time.Now().UnixNano())
}

//endregion

//region PUBLIC METHODS

func (u *Url) Stats() *Stats {
    clicks := repo.GetClicks(u.Id)
    return &Stats{u, clicks}
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

func RegisterClick(id string) {
    repo.RegisterClick(id)
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
