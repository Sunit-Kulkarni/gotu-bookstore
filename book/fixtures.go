package book

import (
	"context"
	_ "embed"
	"encore.app/db"
	"log"

	"encore.dev"
)

//go:embed fixtures.sql
var fixtures string

func init() {
	if encore.Meta().Environment.Cloud == encore.CloudLocal {
		if _, err := db.Bookstoredb.Exec(context.Background(), fixtures); err != nil {
			log.Fatalln("unable to add fixtures:", err)
		}
	}
}
