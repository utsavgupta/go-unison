package migration

import (
	"fmt"

	"cloud.google.com/go/datastore"

	"github.com/utsavgupta/go-unison/example/webapp/ent"
	"github.com/utsavgupta/go-unison/example/webapp/repo"
)

// Apply1584290742 Migrates new albums
func (u *UnisonMigrations) Apply1584290742(t *datastore.Transaction, ns string) error {

	rhcpAlbums := []ent.Album{
		ent.Album{ID: "bytheway", Name: "By the Way"},
		ent.Album{ID: "californication", Name: "Californication"},
	}

	keys := make([]*datastore.Key, len(rhcpAlbums))

	for idx, album := range rhcpAlbums {
		keys[idx] = repo.KeyForAlbum(album.ID, "rhcp", ns)
	}

	_, err := t.PutMulti(keys, rhcpAlbums)

	if err != nil {
		fmt.Println(err)
		return err
	}

	totoAlbums := []ent.Album{
		ent.Album{ID: "iv", Name: "IV"},
		ent.Album{ID: "hydra", Name: "Hydra"},
		ent.Album{ID: "seventhone", Name: "The Seventh One"},
	}

	keys = make([]*datastore.Key, len(totoAlbums))

	for idx, album := range totoAlbums {
		keys[idx] = repo.KeyForAlbum(album.ID, "toto", ns)
	}

	_, err = t.PutMulti(keys, totoAlbums)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
