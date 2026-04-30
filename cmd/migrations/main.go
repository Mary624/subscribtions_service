package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func main() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "subscribes",
	})

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		log.Fatal("failed to run migrations", err.Error())
	}

	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}
