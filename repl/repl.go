package repl

import (
	"bufio"
	"fmt"
	"io"
	"vabna/errs"
	"vabna/evaluator"
	"vabna/lexer"
	"vabna/object"
	"vabna/parser"
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
		evals := evaluator.Eval(prog, env)
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
