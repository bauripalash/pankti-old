package androidapi

import (
	"bytes"
	"io/ioutil"
	"strings"

	"bauri.palash/pankti/evaluator"
	"bauri.palash/pankti/lexer"
	"bauri.palash/pankti/object"
	"bauri.palash/pankti/parser"
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

	env := object.NewEnv()
	eh := evaluator.ErrorHelper{Source: i}
	printBuff := bytes.Buffer{}
	evd := evaluator.Eval(prog, env, eh, &printBuff)
	rd, err := ioutil.ReadAll(&printBuff)

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
