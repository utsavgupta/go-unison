package unison

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	// UnisonMigrationMetaKind is used for storing applied migration data in Datastore
	UnisonMigrationMetaKind = "UnisonMigrationMeta"
)

// The metadata of each unison migrations is stored
// in unisonMigrationMeta
type unisonMigrationMeta struct {
	Name      string    `datastore:"-"`
	Timestamp int64     `datastore:"ts"`
	AppliedAt time.Time `datastore:"applied_at"`
}

// unisonMigrationMetaSet contains a set of unisonMigrationMeta
// The structure has Len, Less, and Swap methods defined
// on them which enabales sorting of the contained unisonMigrationMetas.
type unisonMigrationMetaSet struct {
	set []unisonMigrationMeta
}

// getLastAppliedMigrationTimeStamp returns the UNIX timestamp of the last
// applied migration
func getLastAppliedMigrationTimeStamp(c *datastore.Client, ns string) int64 {
	appliedMigrationsGQL := datastore.NewQuery(UnisonMigrationMetaKind).Namespace(ns).Order("-ts").Limit(1)
	appliedMigrations := []unisonMigrationMeta{}

	_, err := c.GetAll(context.Background(), appliedMigrationsGQL, &appliedMigrations)

	if err != nil {
		panic(err)
	}

	var lastAppliedMigrationTS int64

	if len(appliedMigrations) > 0 {
		lastAppliedMigrationTS = appliedMigrations[0].Timestamp
	}

	return lastAppliedMigrationTS
}

func (ms unisonMigrationMetaSet) Len() int {
	return len(ms.set)
}

func (ms unisonMigrationMetaSet) Less(i, j int) bool {
	return ms.set[i].Timestamp < ms.set[j].Timestamp
}

func (ms unisonMigrationMetaSet) Swap(i, j int) {
	ms.set[i], ms.set[j] = ms.set[j], ms.set[i]
}

func generateMigrationMetaSet(t reflect.Type) (*unisonMigrationMetaSet, error) {

	migrationMetaSet := &unisonMigrationMetaSet{}
	migrationMetaSet.set = make([]unisonMigrationMeta, t.NumMethod())

	migpattern := regexp.MustCompile("^Apply(\\d+)$")
	var groups []string

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		if groups = migpattern.FindStringSubmatch(m.Name); len(groups) != 2 {
			return nil, fmt.Errorf("Invalid migration name %s", m.Name)
		}

		ts, err := strconv.Atoi(groups[1])

		if err != nil {
			return nil, err
		}

		// TODO: validate method signature

		migration := unisonMigrationMeta{Name: groups[0], Timestamp: int64(ts)}

		migrationMetaSet.set[i] = migration
	}

	return migrationMetaSet, nil
}

// RunMigrations ???
func RunMigrations(c *datastore.Client, ns string, migratable interface{}) {

	// Get the timestamp of the last applied migration
	lastAppliedMigrationTS := getLastAppliedMigrationTimeStamp(c, ns)

	tMigratable := reflect.TypeOf(migratable)

	migrationMetaSet, err := generateMigrationMetaSet(tMigratable)

	if err != nil {
		panic(err)
	}

	sort.Sort(migrationMetaSet)

	vMigratable := reflect.ValueOf(migratable)

	for _, migration := range migrationMetaSet.set {

		if migration.Timestamp <= lastAppliedMigrationTS {
			continue
		}

		fmt.Printf("Applying migration %d ... ", migration.Timestamp)

		tx, _ := c.NewTransaction(context.Background())

		method := vMigratable.MethodByName(migration.Name)

		rtx := reflect.ValueOf(tx)
		rns := reflect.ValueOf(ns)

		method.Call([]reflect.Value{rtx, rns})

		txkey := &datastore.Key{Kind: UnisonMigrationMetaKind, Name: migration.Name, Namespace: ns}
		migration.AppliedAt = time.Now()

		tx.Put(txkey, &migration)

		_, err := tx.Commit()

		if err != nil {
			panic(err)
		}

		fmt.Printf("Done\n")
	}

	fmt.Printf("We are done !!\n")
}
