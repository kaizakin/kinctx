package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/kaizakin/kinctx/data"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init kinctx Database",
	Long: `This initializes a sqlite Database with Three tables.
- one for storing static commands
- one for storing templ(dynamic) commands
- this is the table that stores all the aliases that you've created with kinctx.`,
	Run: func(cmd *cobra.Command, args []string) {
		data.CreateTable()
		fmt.Println("Kinctx Datatables created successfully ❇️")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
