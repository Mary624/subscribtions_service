package main

import (
	"fmt"
	"log"
	"subscriptions_rest/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// TODO
	cfg, err := config.New("../../config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.New(
		"file://../../db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.Postgres.Username,
			cfg.Postgres.Password, cfg.Postgres.URL, cfg.Postgres.Database))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
