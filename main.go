package main

import(
	"github.com/kaizakin/kinctx/cmd"
	"github.com/kaizakin/kinctx/data"
) 

func main() {
	data.OpenDatabase() 
	cmd.Execute()
}
