package main

import (
	"fmt"
	"os"

	"github.com/kaizakin/kinctx/cmd"
	"github.com/kaizakin/kinctx/data"
) 

func main() {
	db, err := data.OpenDatabase() 
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close() // making sure db connection is closed after exit

	cmd.Execute()
}
