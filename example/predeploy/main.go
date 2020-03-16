package main

import (
	"context"

	"github.com/utsavgupta/go-unison/example/webapp/migration"

	"cloud.google.com/go/datastore"
	"github.com/utsavgupta/go-unison/unison"
)

const (
	namespace = "unison-demo"
)

func main() {

	dsClient, err := datastore.NewClient(context.Background(), "*detect-project-id*")

	if err != nil {
		panic(err)
	}

	var unisonMigrations migration.UnisonMigrations

	unison.RunMigrations(dsClient, namespace, &unisonMigrations)
}
