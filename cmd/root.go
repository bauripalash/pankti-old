package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{

	Use:   "pankti",
	Short: "Pankti is practical bengali programming language",
	Long:  "Pankti is an interpreted dynamically typed programming language for programming in the Bengali language",

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
