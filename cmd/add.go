/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/kaizakin/kinctx/data"
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
		command, err := captureInputFromEditor()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("parsed command %v", command)

		err = data.AddSnippet(command)
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
