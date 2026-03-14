/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
)


func main(rawCmd string){
	fmt.Println(lipgloss.NewStyle().
	Bold(true).
	Render("\n  ⚡ KIN EXECUTION ENGINE ⚡\n"))

	re := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)(?::=([^}]*))?\}`)
	matches := re.FindAllStringSubmatch(rawCmd, -1)

	if len(matches) == 0 {
		fmt.Println("No placeholders found. Executing immediately...")
		execute(rawCmd)
		return
	}

	answers := make(map[string]*string)
	var fields []huh.Field

	for _, match := range matches {
		fullMatch := match[0]
		key := match[1]
		defaultVal := match[2]

		val := new(string) // huh needs a pointer to mutate values
		*val = defaultVal // assign the default value beforehand

		answers[fullMatch] = val // store the pointer in map

		field := huh.NewInput().
					Title(fmt.Sprintf("Set %s", key)).
					Description(fmt.Sprintf("Default: %s", defaultVal)).
					Value(val)
				
		fields = append(fields, field)
	}

	form := huh.NewForm(
		huh.NewGroup(fields...),
	).WithTheme(huh.ThemeBase16())

	err := form.Run()
	if err != nil {
		fmt.Println("Execution cancelled by user.")
		os.Exit(1)
	}

	//substitute values
	finalCmd := rawCmd
	for fullMatch, answerPtr := range answers {
		userValue := *answerPtr

		if userValue == "" {
			fmt.Printf("Warning: placeholder %s was left empty.\n", fullMatch)
		}

		finalCmd = strings.Replace(finalCmd, fullMatch, userValue, 1)
	}

	execute(finalCmd)
}

func execute(cmdStr string) {
	if strings.TrimSpace(cmdStr) == "" {
		fmt.Println("\nEmpty command cannot be executed!")
		os.Exit(1)
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}

	cmd := exec.Command(shell, "-c", cmdStr)

	//wire up file descriptors
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Printf("\nCommand failed to execute: %v\n", err)
		os.Exit(1)
	}
}

var list []string

func fzf() error {
	for _, s := range SnippetSlice {
		list = append(list, s.Command)
	}

	cmd := exec.Command("fzf", "--height=40%", "--layout=reverse", "--border" , "--read0", "--print0", "--highlight-line")

	cmd.Stdin = strings.NewReader(strings.Join(list, "\x00"))
	cmd.Stderr = os.Stderr

	var out bytes.Buffer // we save the output of the fzf to a buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if exiterror, ok := err.(*exec.ExitError); ok && exiterror.ExitCode() == 130 {
			return fmt.Errorf("Selection cancelled.")
		}

		return fmt.Errorf("Error running fzf: %v\n", err)
	}

	selected := strings.TrimSuffix(out.String(), "\x00")

	if selected == "" {
		return fmt.Errorf("No option selected.")
	}


	err := data.UpdateSnippetUsage(selected)
	if err != nil {
		fmt.Printf("\nWarning: could not update usage count: %v\n", err)
	}

	main(selected)
	return nil
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search, evaluate, and execute a saved command",
	Long: `Use an interactive fzf search to quickly find a saved snippet. Once selected, kinctx will parse any placeholders (like ${VAR:=default}) and prompt you for values using an interactive form. After filling in the variables, the finalized command is directly executed in your terminal. This is the core workflow of kinctx!`,
	Run: func(cmd *cobra.Command, args []string) {
		SnippetSlice, err = data.ListSnippets() // again updating the var becoz what if the user didn't run list for a long time
		// so when user  types search snippetslice var gets updated with latest commands
		if err != nil {
			fmt.Printf("Failed to load commands: %v\n", err)
			return
		}
		err = fzf()
		if err != nil {
			fmt.Printf("Search failed: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
