/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"pankti/ide"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ideCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ideCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
