package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add command to second brain",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add command called")
	},
}

func init() {
	staticCmd.AddCommand(addCmd)
}
