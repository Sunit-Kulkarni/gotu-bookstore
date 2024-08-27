package db

import (
	"encore.dev/storage/sqldb"
)

// Create the bookstore database and assign it to the "bookstoredb" variable
var Bookstoredb = sqldb.NewDatabase("bookstore", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
