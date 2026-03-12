package data

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite" // just import this for registering the driver not actually using it so _ in front.
	sqlitelib "modernc.org/sqlite/lib"
)

var db *sql.DB

type Snippets struct {
	Id int
	Command string
	UsageCount int
	createdAt time.Time
	CreatedAtFormatted string
}

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

func ListSnippets() ([]Snippets, error) {
	query := `SELECT id, command, usage_count, created_at from snippets
	ORDER BY id ASC`

	listCommandstmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("List query preparation failed!")
	}
	defer listCommandstmt.Close()

	rows, err := listCommandstmt.Query() // Exec is for writing ops and query is for read ops
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var snippetSlice []Snippets

	for rows.Next() {
		var s Snippets

		err := rows.Scan(&s.Id, &s.Command, &s.UsageCount, &s.createdAt)
		if err != nil {
			return nil, err
		}

		s.CreatedAtFormatted = s.createdAt.Format("2006-01-02")

		snippetSlice = append(snippetSlice, s)
	}

	return snippetSlice, err
}

func DeleteSnippets(idsToDelete []int) error {
	if len(idsToDelete) == 0 {
		return fmt.Errorf("No id's to Delete!")
	}

	// innorder to prevent sql injection and as a safe good practice construct a placeholder string and pass
	// args as an arguement
	placeHolders := make([]string, len(idsToDelete))
	args := make([]any, len(idsToDelete))

	for i, id := range idsToDelete{
		placeHolders[i] = "?"
		args[i] = id
	}

	placeholderStr := strings.Join(placeHolders, ",")
	query := fmt.Sprintf("DELETE FROM snippets WHERE id IN (%s)", placeholderStr)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("Failed to delete commands: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Deleted but could not verify dropped rows: %w", err)
	}

	fmt.Printf("%d rows affected\n", rowsAffected)

	return nil
}