package androidapi

import (
	"bytes"
	"io"
	"strings"

	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

func NewLexer(l string) lexer.Lexer {
	return lexer.NewLexer(l)
}

func NewParser(l lexer.Lexer) parser.Parser {
	return *parser.NewParser(&l)
}

func DoParse(i string) string {
	l := lexer.NewLexer(i)
	p := parser.NewParser(&l)
	prog := p.ParseProg()

	if len(p.GetErrors()) >= 1 {
		tempErrs := []string{}

		for _, item := range p.GetErrors() {
			tempErrs = append(tempErrs, item.String())
		}

		return strings.Join(tempErrs, " \n")
	}

	env := object.NewEnvMap()
	eh := object.ErrorHelper{Source: i}
	printBuff := bytes.Buffer{}
	evd := evaluator.Eval(prog, env, eh, &printBuff, true)
	rd, err := io.ReadAll(&printBuff)

	printValue := ""

	if err == nil {
		printValue = string(rd[:])
	}

	if evd != nil {
		return printValue + evd.Inspect()
	} else {
		return printValue
	}
}
