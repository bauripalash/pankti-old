//go:build !noide
// +build !noide

package cmd

import (
	"github.com/spf13/cobra"
	"go.cs.palashbauri.in/pankti/ide"
)

// ideCmd represents the ide command
var ideCmd = &cobra.Command{
	Use:   "ide",
	Short: "Run GUI Editor",
	Long:  `Pankti IDE is basic editor to quickly write and run Pankti Programs`,
	Run: func(cmd *cobra.Command, args []string) {
		ide.RunIde()
	},
}

func init() {
	rootCmd.AddCommand(ideCmd)
}
