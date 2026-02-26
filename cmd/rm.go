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
	"strconv"
	"strings"

	"github.com/kaizakin/kinctx/data"
	"github.com/spf13/cobra"
)

var fzfInput strings.Builder

func rm() error {

	for _, s := range SnippetSlice {
		fzfInput.WriteString(fmt.Sprintf("%d\t%s\x00", s.Id, s.Command)) // null byte is the delimiter
	}

	cmd := exec.Command("fzf",
		"--height=40%", "--layout=reverse",
		"--border", "--read0", "--highlight-line", "--multi",
		"--delimiter=\\t", "--with-nth=2..",
		"--print0",
		"--header=Use [TAB] to select, [ENTER] to delete")

	cmd.Stdin = strings.NewReader(fzfInput.String())
	cmd.Stderr = os.Stderr
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 130 {
			fmt.Println("CANCELLED...")
			return nil
		}

		return fmt.Errorf("Fzf error %v", err)
	}

	// With --print0, each selected entry is NUL-separated, so multiline
	// commands stay as a single record instead of being split on '\n'.
	rawOut := strings.TrimRight(out.String(), "\x00") // first remove the trailine null
	selectedEntries := strings.Split(rawOut, "\x00")

	if len(selectedEntries) == 0 || (len(selectedEntries) == 1 && strings.TrimSpace(selectedEntries[0]) == "") {
		return fmt.Errorf("No selected Commands to delete")
	}

	var idsToDelete []int

	for _, entry := range selectedEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Expect "ID<TAB>command..."; command itself may contain newlines.
		parts := strings.SplitN(entry, "\t", 2)
		if len(parts) < 2 {
			continue
		}

		idStr := strings.TrimSpace(parts[0])
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		idsToDelete = append(idsToDelete, id)
	}

	if len(idsToDelete) == 0 {
		return fmt.Errorf("No valid command IDs selected to delete")
	}

	err := data.DeleteSnippets(idsToDelete)
	if err != nil {
		return err
	}

	return nil
}

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
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

		var err error
		err = rm();
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
