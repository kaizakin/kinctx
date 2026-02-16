/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
)

func getEditor() string{
	editor := os.Getenv("EDITOR")
	if editor == ""{
		editor = "nano"
	}
	
	fmt.Printf("found %v\n", editor)

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

			fmt.Printf("command read from stdin %s \n", inputCommand)
		}else{
			inputCommand, err = captureInputFromEditor()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("parsed command %v", inputCommand)
		}

		err = data.AddSnippet(inputCommand)
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("command saved successfully")	
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
