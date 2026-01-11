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

	snippetTableSQL := `
	CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,          
    command TEXT NOT NULL,       
    usage_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	aliasTableSQL := `
	CREATE TABLE IF NOT EXISTS aliases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias_name TEXT UNIQUE NOT NULL, 
    command TEXT NOT NULL    
	);`

	createSnippetTableStatement, err := db.Prepare(snippetTableSQL)
    if err != nil {
        log.Fatal(err.Error())
    }

	createAliasTableStatement, err := db.Prepare(aliasTableSQL)
	if err != nil{
		log.Fatal(err.Error())
	}

	createSnippetTableStatement.Exec()
	log.Println("Snippets table created")

	createAliasTableStatement.Exec()
	log.Println("Alias Table Created")
}