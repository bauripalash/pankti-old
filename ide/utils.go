//go:build !noide
// +build !noide

package ide

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

func OpenFile(filename string) (string, error) {

	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SaveFile(filename string, content string, overwrite bool) error {
	if overwrite {
		err := os.WriteFile(filename, []byte(content), 0644)
		return err
	} else {
		nf, err := os.Create(filename)
		defer nf.Close()
		if err != nil {
			return err
		}

		nf.Write([]byte(content))
		nf.Sync()
	}

	return nil
}

func RunFile(src string) string {

	l := lexer.NewLexer(src)
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
	eh := object.ErrorHelper{Source: src}
	printBuff := bytes.Buffer{}
	evd := evaluator.Eval(prog, env, eh, &printBuff, true)
	//rd, _ := ioutil.ReadAll(&printBuff)
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
