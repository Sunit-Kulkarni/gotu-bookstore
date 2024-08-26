package account

import (
	"encore.dev/storage/sqldb"
)

// Create the bookstore database and assign it to the "bookstoredb" variable
var bookstoredb = sqldb.NewDatabase("bookstore", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
