/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
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
		fmt.Println("No command was provided. Please pipe a command or type one in the editor.")
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
	Short: "Save a new command to your kinctx store",
	Long: `Securely store a shell command in your kinctx database. You can either pipe a command via stdin or open your default editor ($EDITOR) to compose complex, multi-line commands. You can also include dynamic placeholders like ${VAR:=default} to be evaluated at runtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputCommand string
		var err error

		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Printf("Error reading stdin stats: %v\n", err)
			return
		}

		// if modechardevice bit is not set means data is being piped in
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// read everything from the pipe
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("Error reading from pipe: %v\n", err)
				return
			}

			inputCommand = strings.TrimSpace(string(bytes))

		}else{
			inputCommand, err = captureInputFromEditor()
			if err != nil {
				fmt.Printf("Error capturing input from editor: %v\n", err)
				return
			}

		}

		err = ValidateCmd(inputCommand)
		if err != nil {
			fmt.Printf("\nCommand Validation failed: %v\n", err)
			return
		}

		err = data.AddSnippet(inputCommand)
		if err != nil {
			fmt.Printf("Failed to save command: %v\n", err)
		}else{
			fmt.Println("Command saved successfully!")	
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
