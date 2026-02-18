/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/shlex"
	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
)


func main(rawCmd string){
	fmt.Println(lipgloss.NewStyle().
	Foreground(lipgloss.Color("#bd93f9")).
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
	).WithTheme(huh.ThemeDracula())

	err := form.Run()
	if err != nil {
		fmt.Println("User cancelled!")
		os.Exit(1)
	}

	//substitute values
	finalCmd := rawCmd
	for fullMatch, answerPtr := range answers {
		userValue := *answerPtr

		if userValue == "" {
			fmt.Printf("Warning: %s was left empty.\n", fullMatch)
		}

		finalCmd = strings.Replace(finalCmd, fullMatch, userValue, 1)
	}

	execute(finalCmd)
}

func execute(cmdStr string) {
	fmt.Printf("Executing %s", cmdStr)

	parts, err := shlex.Split(cmdStr) 
	if err != nil {
		fmt.Println("Failed to parse the command!")
		os.Exit(1)
	}

	if len(parts) == 0 {
		fmt.Println("Empty command cannot be executed!")
		os.Exit(1)
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	//wire up file descriptors
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("Command failed to execute!")
		os.Exit(1)
	}
}

var list []string

func fzf() error {
	for _, s := range SnippetSlice {
		list = append(list, s.Command)
	}

	cmd := exec.Command("fzf", "--height=40%", "--layout=reverse", "--border")

	cmd.Stdin = strings.NewReader(strings.Join(list, "\n"))
	cmd.Stderr = os.Stderr

	var out bytes.Buffer // we save the output of the fzf to a buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if exiterror, ok := err.(*exec.ExitError); ok && exiterror.ExitCode() == 130 {
			return fmt.Errorf("Selection cancelled.")
		}

		return fmt.Errorf("Error running fzf: %v\n", err)
	}

	selected := strings.TrimSpace(out.String())

	if selected == "" {
		return fmt.Errorf("No option selected.")
	}

	fmt.Printf("picked command %s", selected)
	main(selected)
	return nil
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		SnippetSlice, err = data.ListSnippets() // again updating the var becoz what if the user didn't run list for a long time
		// so when user  types search snippetslice var gets updated with latest commands
		if err != nil {
			log.Fatal(err)
		}
		err = fzf()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("search and execute successfull")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
