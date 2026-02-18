/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/shlex"
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
		main("docker run -p ${HOST_PORT:=5432}:5432 -e PASS=${DB_PASS:=secret} postgres:${PG_VERSION:=16}")
		fmt.Println("ended successfully!")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
