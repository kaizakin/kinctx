/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
	"github.com/charmbracelet/lipgloss"
)

var snippetSlice []data.Snippets
var err error

var cardStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true). // Left border only
	BorderForeground(lipgloss.Color("#5B4FB5")).                // Purple border
	PaddingLeft(2).
	MarginBottom(1)

var commandStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF79C6")). // Pink/Magenta 
	Bold(true)

var metaStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#6272A4")) // Muted Blue/Grey

var statsStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#50FA7B")) // Bright Green

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		snippetSlice, err = data.ListSnippets()
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range snippetSlice {
			renderSnippet(s)
		}
	},
}
func renderSnippet(s data.Snippets) {
	cmdBlock := commandStyle.Render(fmt.Sprintf("$ %s", s.Command))

	idStr := fmt.Sprintf("#%d", s.Id)
	dateStr := s.CreatedAtFormatted
	
	usageStr := statsStyle.Render(fmt.Sprintf("Used %d times", s.UsageCount))
	
	metaBlock := metaStyle.Render(strings.Join([]string{idStr, dateStr, usageStr}, " • "))

	ui := cardStyle.Render(fmt.Sprintf("%s\n%s", cmdBlock, metaBlock))

	fmt.Println(ui)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
