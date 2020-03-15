package repo

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/utsavgupta/go-unison/example/webapp/ent"
)

const (
	artistKind = "Artist"
)

type ArtistLoader interface {
	LoadAllArtists(context.Context) ([]ent.Artist, error)
}

func NewArtistLoader(client *datastore.Client, ns string) ArtistLoader {
	return &artistLoader{dsClient: client, dsNamespace: ns}
}

type artistLoader struct {
	dsClient    *datastore.Client
	dsNamespace string
}

func (al *artistLoader) LoadAllArtists(ctx context.Context) ([]ent.Artist, error) {

	var artists []ent.Artist

	_, err := al.dsClient.GetAll(ctx, datastore.NewQuery(artistKind).Namespace(al.dsNamespace), &artists)

	return artists, err
}

func KeyForArtist(artistID string, ns string) *datastore.Key {
	key := datastore.NameKey(artistKind, artistID, nil)
	key.Namespace = ns

	return key
}
