package cmd

import (
	"os"
	"pankti/repl"

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// replCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// replCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
