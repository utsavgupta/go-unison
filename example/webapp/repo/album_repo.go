package repo

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/utsavgupta/go-unison/example/webapp/ent"
)

const (
	albumKind = "Album"
)

type AlbumLoader interface {
	LoadAllAlbums(context.Context) ([]ent.Album, error)
	LoadAllAlbumsByArtist(context.Context, string) ([]ent.Album, error)
}

func NewAlbumLoader(client *datastore.Client, ns string) AlbumLoader {
	return &albumLoader{dsClient: client, dsNamespace: ns}
}

type albumLoader struct {
	dsClient    *datastore.Client
	dsNamespace string
}

func (al *albumLoader) LoadAllAlbums(ctx context.Context) ([]ent.Album, error) {

	var albums []ent.Album

	_, err := al.dsClient.GetAll(ctx, datastore.NewQuery(albumKind).Namespace(al.dsNamespace), &albums)

	return albums, err
}

func (al *albumLoader) LoadAllAlbumsByArtist(ctx context.Context, artistID string) ([]ent.Album, error) {

	var albums []ent.Album

	_, err := al.dsClient.GetAll(ctx, datastore.NewQuery(albumKind).Namespace(al.dsNamespace).Ancestor(KeyForArtist(artistID, al.dsNamespace)), &albums)

	return albums, err
}

func KeyForAlbum(albumID string, artistID string, ns string) *datastore.Key {
	key := datastore.NameKey(albumKind, albumID, KeyForArtist(artistID, ns))
	key.Namespace = ns

	return key
}
