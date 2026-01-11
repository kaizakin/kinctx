package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kinctx",
	Short: "A minimalist second brain for your terminal snippets",
	Long: `kinctx is a high-performance CLI snippet manager designed for 
developers who live in the terminal. 

It allows you to:
- Save complex commands you keep on forgetting.
- Save commands with dynamic placeholders ({{VAR}}) and dynamically fill values at execution time.
- Kinctx also acts as a alias manager.`,
}

func Execute() {
	err := rootCmd.Execute() 
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "0.1.0"

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


