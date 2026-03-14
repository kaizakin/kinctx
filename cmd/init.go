package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/kaizakin/kinctx/data"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the local kinctx database",
	Long: `Creates the necessary SQLite database and tables in your local environment. This is typically run once when setting up kinctx for the first time. It sets up storage for static commands, dynamic templated commands, and alias mappings.`,
	Run: func(cmd *cobra.Command, args []string) {
		data.CreateTable()
		fmt.Println(`
 _  ___            _        
| |/ (_)_ __   ___| |___  __
| ' /| | '_ \ / __| __\ \/ /
| . \| | | | | (__| |_ >  < 
|_|\_\_|_| |_|\___|\__/_/\_\
`)
		fmt.Println("Kinctx Datatables created successfully ❇️")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
