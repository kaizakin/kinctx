package main

import (
	"log"

	"github.com/kaizakin/kinctx/cmd"
	"github.com/kaizakin/kinctx/data"
) 

func main() {
	db, err := data.OpenDatabase() 
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // making sure db connection is closed after exit

	cmd.Execute()
}
