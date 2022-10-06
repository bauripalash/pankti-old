package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"pankti/evaluator"
	"pankti/lexer"
	"pankti/object"
	"pankti/parser"
	"pankti/repl"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [FILENAME]",
	Short: "Run a Pankti Source File",
	Long:  `Run a pankti source file providing as a argument to this command`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please provile a file to run")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			filename := args[0]
			_, err := os.Stat(filename)

			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("File `%s` does not exist!\n\n", filename)
			}

			f, err := os.ReadFile(filename)

			if err != nil {
				fmt.Printf("Cannot read `%s`\n\n", filename)
			}

			//fmt.Println(string(f))

			lx := lexer.NewLexer(string(f))
			ps := parser.NewParser(&lx)
			at := ps.ParseProg()

			if len(ps.GetErrors()) != 0 {
				repl.ShowParseErrors(os.Stdin, ps.GetErrors())
				fmt.Printf("fix above mentioned errors first!\n\n")
			}
			env := object.NewEnv()
			evd := evaluator.Eval(at, env)

			if evd != nil {
				fmt.Println(evd.Inspect())
			}

			//fmt.Println(args[0])

		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
