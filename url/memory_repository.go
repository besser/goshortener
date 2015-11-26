package url

//region TYPES

type repositoryMemory struct {
	urls map[string]*Url
}

//endregion

//region PUBLIC FUNCTIONS

func (r *repositoryMemory) IdExists(id string) bool {
	_, exist := r.urls[id]
	return exist
}

func (r *repositoryMemory) FindById(id string) *Url {
	return r.urls[id]
}

func (r *repositoryMemory) FindByUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Destination == url {
			return u
		}
	}

	return nil
}

func NewRepoMem() *repositoryMemory {
	return &repositoryMemory{make(map[string]*Url)}
}

func (r *repositoryMemory) Save(url Url) error {
	r.urls[url.Id] = &url
	return nil
}

//endregion