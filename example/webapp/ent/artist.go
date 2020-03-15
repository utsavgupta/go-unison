package ent

type Artist struct {
	ID   string `json:"id" datastore:"id"`
	Name string `json:"name" datastore:"name"`
}

type Artists struct {
	Items []Artist `json:"items,omitempty"`
}
