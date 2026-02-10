package data

import (
	"database/sql"
	"fmt"
	"log"

	"modernc.org/sqlite"
	// _ "modernc.org/sqlite" // just import this for registering the driver not actually using it so _ in front.
	sqlitelib "modernc.org/sqlite/lib"
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
    command TEXT UNIQUE NOT NULL,       
    usage_count INTEGER DEFAULT 0,
	last_used_at INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createSnippetTableStatement, err := db.Prepare(snippetTableSQL)
	
    if err != nil {
        if sqliteErr, ok := err.(*sqlite.Error); ok {
			if(sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE){
				fmt.Println("You already have this command saved Einstein, try searching for it!")
			}
		}else{
			log.Fatal(err)
		}
    }

	createSnippetTableStatement.Exec()
	log.Println("Snippets table created")
}