package data

import (
	"database/sql"
	"log"
	"fmt"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite" // just import this for registering the driver not actually using it so _ in front.
	sqlitelib "modernc.org/sqlite/lib"
)

var db *sql.DB

func OpenDatabase() (*sql.DB, error) {
	var err error

	db, err = sql.Open("sqlite", "./sqlite-database.db")
	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTable(){

	snippetTableSQL := `
	CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    command TEXT UNIQUE NOT NULL,       
    usage_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(snippetTableSQL)
	
	if err != nil {
		log.Fatal(err)
	}


	log.Println("Snippets table created")
}

func AddSnippet(command string) error {
	query := `INSERT INTO snippets (command) VALUES (?)`

	insertCommandStmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer insertCommandStmt.Close()

	_, err = insertCommandStmt.Exec(command)
	if err != nil {
        if sqliteErr, ok := err.(*sqlite.Error); ok {
			if(sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE){
				return fmt.Errorf("You already have this command saved Einstein, try searching for it!")
			}
		}else{
			return err
		}
    }
	
	return nil
}
