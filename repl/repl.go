package repl

import (
	"bufio"
	"fmt"
	"io"
	"bauri.palash/pankti/errs"
	"bauri.palash/pankti/evaluator"
	"bauri.palash/pankti/lexer"
	"bauri.palash/pankti/object"
	"bauri.palash/pankti/parser"
)

const PROMPT = "-> "

func Repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

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
		eh := evaluator.ErrorHelper{Source: input}
		evals := evaluator.Eval(prog, env, eh)
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
