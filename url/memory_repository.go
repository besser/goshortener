package url

//region TYPES

type repositoryMemory struct {
	clicks map[string]int
	urls   map[string]*Url
}

//endregion

//region PUBLIC METHODS

func (r *repositoryMemory) GetClicks(id string) int  {
    return r.clicks[id]
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

func (r *repositoryMemory) IdExists(id string) bool {
	_, exist := r.urls[id]
	return exist
}

func (r *repositoryMemory) RegisterClick(id string) {
    r.clicks[id] += 1
}

func (r *repositoryMemory) Save(url Url) error {
	r.urls[url.Id] = &url
	return nil
}

//endregion

//region PUBLIC FUNCIONS

func NewRepoMem() *repositoryMemory {
	return &repositoryMemory{
		make(map[string]int),
		make(map[string]*Url),
	}
}

//endregion