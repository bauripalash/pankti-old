package main

// #import <stdlib.h>
import "C"

import (
	"bytes"
	"io/ioutil"
	"strings"
    "os"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

func init() {
    
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}


//export DoParse
func DoParse(src string) *C.char{
    
    lx := lexer.NewLexer(src)
    p := parser.NewParser(&lx)
    
    prog := p.ParseProg()
    
    if len(p.GetErrors()) >= 1 {
		tempErrs := []string{}

		for _, item := range p.GetErrors() {
			tempErrs = append(tempErrs, item.String())
		}

		return C.CString(strings.Join(tempErrs, " \n"))
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
		return C.CString(printValue + evd.Inspect())
	} else {
		return C.CString(printValue)
	}
}

func main(){}
