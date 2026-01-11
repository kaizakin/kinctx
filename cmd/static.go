package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "Static snippet store",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("static command called")
	},
}

func init() {
	rootCmd.AddCommand(staticCmd)
}
