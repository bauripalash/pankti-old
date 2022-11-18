package cmd

import (
	"os"

	"go.cs.palashbauri.in/pankti/repl"

	"github.com/spf13/cobra"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Run a quick REPL (not recommended)",
	Long:  `A quick REPL inside current Terminal (not recommended; use IDE instead)`,
	Run: func(cmd *cobra.Command, args []string) {
		repl.Repl(os.Stdin, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(replCmd)

}
