package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"syscall/js"

	"github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
}

func main() {
	js.Global().Set("runner", runner())
	js.Global().Set("evs", js.FuncOf(rfile))
	<-make(chan bool)
}

//export DoRun
func DoRun(src string) string {

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
	evd := evaluator.Eval(prog, env, eh, &printBuff, false)
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

func rfile(this js.Value, args []js.Value) interface{} {
	if len(args) >= 1 {
		source_code := args[0].String()
		return DoRun(source_code)
	}
	return ""
}

//export update
func runner() js.Func {
	runFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid arguments"
		}

		source_code := args[0].String()

		output := DoRun(source_code)

		return output

	})

	return runFunc
}
