//go:build !noide
// +build !noide

package ide

import (
	"os"
	"bauri.palash/pankti/evaluator"
	"bauri.palash/pankti/lexer"
	"bauri.palash/pankti/object"
	"bauri.palash/pankti/parser"
	"strings"
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

	env := object.NewEnv()
	eh := evaluator.ErrorHelper{Source: src}
	evd := evaluator.Eval(prog, env, eh)

	if evd != nil {
		return evd.Inspect()
	} else {
		return ""
	}

}
