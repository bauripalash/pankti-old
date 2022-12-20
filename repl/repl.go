package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

const PROMPT = "-> "

// Deprecated
func Repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvMap()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		input := scanner.Text()
		rlexer := lexer.NewLexer(input)
		//fmt.Println(rlexer.)
		/*
		   for !rlexer.AtEOF(){
		       fmt.Println(rlexer.NextToken().Type)
		   }
		*/
		p := parser.NewParser(&rlexer)

		prog := p.ParseProg()
		fmt.Println(prog.Stmts)

		if len(p.GetErrors()) != 0 {
			ShowParseErrors(out, p.GetErrors())
			continue
		}
		eh := object.ErrorHelper{Source: input}

		printBuff := bytes.Buffer{}
		evals := evaluator.Eval(prog, env, eh, &printBuff, false)

		if evals != nil {
			//fmt.Println(evals)
			io.WriteString(out, evals.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func ShowParseErrors(out io.Writer, errs []errs.ParserError) {
	for _, msg := range errs {
		io.WriteString(out, "\t ERR >"+msg.String()+"\n")
	}
}
