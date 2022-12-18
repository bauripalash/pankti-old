package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"os"

	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
	"go.cs.palashbauri.in/pankti/repl"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [FILENAME]",
	Short: "Run a Pankti Source File",
	Long:  `Run a pankti source file providing as a argument to this command`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a file to run")
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
			/*
			   for !lx.AtEOF(){
			       nt := lx.NextToken()
			       fmt.Printf(` Lit-> %s
			                   Line -> %d
			                   Col -> %d
			                   Type -> %s

			       ` + "\n" , nt.Literal , nt.LineNo , nt.Column , nt.Type)
			   }*/
			ps := parser.NewParser(&lx)
			at := ps.ParseProg()

			if len(ps.GetErrors()) != 0 {
				repl.ShowParseErrors(os.Stdin, ps.GetErrors())
				fmt.Printf("fix above mentioned errors first!\n\n")
			} else {
				env := object.NewEnvMap()
				eh := evaluator.ErrorHelper{Source: string(f)}
				printBuff := bytes.Buffer{}

				evd := evaluator.Eval(at, env, eh, &printBuff, false)

				//rd, _ := ioutil.ReadAll(&printBuff)
				rd, err := ioutil.ReadAll(&printBuff)

				printValue := ""

				if err != nil {
					printValue = string(rd[:])
				}
				fmt.Println(printValue)

				if evd != nil {
					fmt.Println(evd.Inspect())
				}
			}

			//fmt.Println(args[0])

		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

}
