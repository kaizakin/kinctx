package data

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // just import this for registering the driver not actually using it so _ in front.
)

var db *sql.DB

func OpenDatabase() error {
	var err error

	db, err = sql.Open("sqlite", "./sqlite-database.db")
	if err != nil{
		return err
	}

	return db.Ping()
}

func CreateTable(){
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS snippets(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	command TEXT,
	description TEXT,
	tags TEXT,
	usage_count INTEGER DEFAULT 0,
	last_used DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	statement, err := db.Prepare(createTableSQL) // altho it doesn't make any sense to prepare the query for a cli app
	// just for the good practice. prepare actually caches the query to perform efficienct reuse of query and some security from prompt injections.

	if err != nil{
		log.Fatal(err.Error())
	}

	statement.Exec()
	log.Println("snippets table created")
}