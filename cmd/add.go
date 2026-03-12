/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
)

func getEditor() string{
	editor := os.Getenv("EDITOR")
	if editor == ""{
		editor = "nano"
	}

	return editor
}

func captureInputFromEditor() (string, error) {
	tempFile, err := os.CreateTemp("", "kin-*.txt") // create a temp file in the default temp folder
	if err != nil {
		return "", err
	}

	defer os.Remove(tempFile.Name())

	tempFile.Close() // close it, else editor won't open opened files.

	editor := getEditor()

	cmd := exec.Command(editor, tempFile.Name())
	// give em file descriptors to the editor process
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(tempFile.Name())
	parsedCommand := string(bytes)

	if err != nil {
		return "", nil
	}	

	return parsedCommand, nil
}

var placeholderRe = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)(?::=([^}]*))?\}`)
var brokenSyntaxRe = regexp.MustCompile(`\$\{[^}]*$`)

func ValidateCmd(cmd string) error {
	// handle empty or whitespace-only commands
	if strings.TrimSpace(cmd) == "" {
		fmt.Println("no command was found gang!")
		os.Exit(0)
	}

	if brokenSyntaxRe.MatchString(cmd) {
		return errors.New("syntax error: detected unclosed '${' or malformed placeholder 🥲")
	} /// early return if the regex pattern in cmd is broken

	if strings.Contains(cmd, "${}") || strings.Contains(cmd, "${:=}") {
		return errors.New("syntax error: placeholder key cannot be empty 🥲")
	}

	matches := placeholderRe.FindAllStringSubmatch(cmd, -1)
	if len(matches) == 0 {
		return nil
	}// regular command no placeholder syntax

	return nil
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputCommand string
		var err error

		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Println("error reading stdin stats: %w", err)
			return
		}

		// if modechardevice bit is not set means data is being piped in
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// read everything from the pipe
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("error reading from pipe: %w", err)
				return
			}

			inputCommand = strings.TrimSpace(string(bytes))

		}else{
			inputCommand, err = captureInputFromEditor()
			if err != nil {
				log.Fatal(err)
			}

		}

		err = ValidateCmd(inputCommand)
		if err != nil {
			fmt.Println("")
			fmt.Println(err)
			return
		}

		err = data.AddSnippet(inputCommand)
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("command saved successfully!")	
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
