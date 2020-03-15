package migration

import (
	"github.com/utsavgupta/go-unison/example/webapp/ent"
	"github.com/utsavgupta/go-unison/example/webapp/repo"

	"cloud.google.com/go/datastore"
)

// Apply1584290710 Migrates new artists
func (u *UnisonMigrations) Apply1584290710(t *datastore.Transaction, ns string) error {

	artists := []ent.Artist{
		ent.Artist{ID: "whitesnake", Name: "Whitesnake"},
		ent.Artist{ID: "rhcp", Name: "Red Hot Chili Peppers"},
		ent.Artist{ID: "toto", Name: "Toto"},
	}

	keys := make([]*datastore.Key, len(artists))

	for idx, artist := range artists {
		keys[idx] = repo.KeyForArtist(artist.ID, ns)
	}

	_, err := t.PutMulti(keys, artists)

	return err
}
