package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"syscall/js"

	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)


func main(){
	js.Global().Set("runner" , runner())
	<- make(chan bool)
}

//export DoRun
func DoRun(src string) string{
	
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
	printBuff := bytes.Buffer{}
	evd := evaluator.Eval(prog, env, eh, &printBuff, true)
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


//export update
func runner() js.Func{
	runFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1{
			return "Invalid arguments"
		}

		source_code := args[0].String()
		
		output := DoRun(source_code)

		return output

	})

	return runFunc
}
