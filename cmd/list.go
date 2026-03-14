/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
)

var SnippetSlice []data.Snippets
var err error

var gappedBorder = lipgloss.Border{
	Top:        "",
	Bottom:     "─",
	Left:       "|",
	Right:      "|",
	TopLeft:    "",
	TopRight:   "",
	BottomLeft: "",
	BottomRight:"",
}


var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Padding(0, 1)

	commandStyle = lipgloss.NewStyle().
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("110"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")) // Muted separator color

	placeholderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("221")). // Custom color for placeholders
			Italic(true)
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved commands in a styled table",
	Long: `Displays all the snippets you've saved in your kinctx store. The list is presented in a readable, bordered table showing the command itself, its usage count, and when it was created. Placeholders within commands are highlighted for easy readability.`,
	Run: func(cmd *cobra.Command, args []string) {
		SnippetSlice, err = data.ListSnippets()
		if err != nil {
			fmt.Printf("Failed to list snippets: %v\n", err)
			return
		}

		if len(SnippetSlice) == 0 {
			fmt.Println("Your kinctx store is empty. Use 'kinctx add' to save some snippets!")
			return
		}

		renderTable(SnippetSlice)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func formatStyledCommand(cmd string) string {
	re := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)(?::=([^}]*))?\}`)

	lines := strings.Split(cmd, "\n")
	for i, line := range lines {
		matches := re.FindAllStringIndex(line, -1)

		if len(matches) == 0 {
			lines[i] = commandStyle.Render(line)
			continue
		}

		var result string
		var lastEnd int

		for _, match := range matches {
			start, end := match[0], match[1]

			if start > lastEnd {
				result += commandStyle.Render(line[lastEnd:start])
			}

			result += placeholderStyle.Render(line[start:end])
			lastEnd = end
		}

		if lastEnd < len(line) {
			result += commandStyle.Render(line[lastEnd:])
		}
		
		lines[i] = result
	}

	return lipgloss.NewStyle().Align(lipgloss.Left).Render(strings.Join(lines, "\n"))
}

func renderTable(snippets []data.Snippets) {
	rows := [][]string{}

	for _, s := range snippets {
		rows = append(rows, []string{
			formatStyledCommand(s.Command),
			successStyle.Render(fmt.Sprintf("%d uses", s.UsageCount)),
			dimStyle.Render(s.CreatedAtFormatted),
		})
	}

	t := table.New().
		Border(gappedBorder).
		BorderStyle(borderStyle).
		Headers("COMMAND", "USAGE", "DATE").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == -1 {
				return headerStyle
			}
			return lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Left)
		})

	fmt.Println(t.Render())
}