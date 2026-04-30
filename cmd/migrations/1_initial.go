package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	// TODO
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table subscriptions...")
		_, err := db.Exec(`CREATE TABLE subscriptions(
    	user_id uuid NOT NULL,
    	service_name VARCHAR(50) NOT NULL,
    	price integer NOT NULL,
		
);`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table my_table...")
		_, err := db.Exec(`DROP TABLE my_table`)
		return err
	})
}
