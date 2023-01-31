package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"go.cs.palashbauri.in/pankti/compiler"
	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
	"go.cs.palashbauri.in/pankti/vm"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var input = `dhori fib = ekti kaj(x)
jodi (x == 0) tahole 
	ferao(0)
nahole jodi (x==1) tahole
		ferao(1)
	nahole 
		fib(x-1) + fib(x-2)

	sesh

sesh
sesh

fib(10)`

func main() {
	flag.Parse()
	var duration time.Duration
	var result object.Obj

	l := lexer.NewLexer(input)
	p := parser.NewParser(&l)
	prog := p.ParseProg()

	if *engine == "vm" {
		comp := compiler.NewCompiler()
		if err := comp.Compile(prog); err != nil {
			fmt.Printf("compiler error %s", err)
			return
		}
		os.WriteFile("fib.pankc", comp.ByteCode().Instructions, 0644)
		return
		mc := vm.NewVM(*comp.ByteCode())
		start := time.Now()
		if err := mc.Run(); err != nil {
			fmt.Printf("vm error :%s", err)
			return
		}

		duration = time.Since(start)
		result = mc.LastPoppedStackItem()
	} else {
		env := object.NewEnvMap()
		eh := object.ErrorHelper{Source: input}
		start := time.Now()
		result = evaluator.Eval(prog, env, eh, &bytes.Buffer{}, false)
		duration = time.Since(start)

	}

	fmt.Printf("engine=%s; result=%s; duration=%s\n", *engine, result.Inspect(), duration)
}
