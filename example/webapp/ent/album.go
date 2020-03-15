package ent

type Album struct {
	ID   string `json:"id" datastore:"id"`
	Name string `json:"name" datastore:"name"`
}

type Albums struct {
	Items []Album `json:"items,omitempty"`
}
