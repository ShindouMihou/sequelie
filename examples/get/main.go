package main

import (
	"database/sql"
	"log"
	"sequelie"
)

func main() {
	if err := sequelie.ReadDirectory("examples/"); err != nil {
		log.Fatal("failed to read dirs: ", err)
	}
	// For example purposes, we'll use this. You are free to use anything as long as it can support
	// using raw queries since that is the core of what Sequelie does.
	sequel, err := sql.Open("postgres", "postgres://127.0.0.1:5432")
	if err != nil {
		log.Fatal("failed to connect to postgres: ", err)
	}
	// In this example, we are retrieving the "books.get" query from the books.sql and uses
	// the database/sql lib to handle the insertion of the id's value.
	rows, err := sequel.Query(sequelie.Get("books.get"), 0)
	if err != nil {
		log.Fatal("failed to get rows in postgres: ", err)
		return
	}
	defer rows.Close()
	// You can then do whatever you'd like to do with the rows, if you prefer doing it manually.
}
