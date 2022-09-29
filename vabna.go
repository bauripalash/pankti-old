package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	/*
		"vabna/evaluator"
		"vabna/lexer"
		"vabna/object"
		"vabna/parser"
	*/
	"vabna/evaluator"
	"vabna/lexer"
	"vabna/object"
	"vabna/parser"
	"vabna/repl"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}

func main() {
	/*
	   	examplecode := `
	           let newAdder = fn(x) { fn(y) {x+y} };
	           let addTwo = newAdder(2);
	           addTwo

	       `

	   	l := lexer.NewLexer(examplecode)
	   	p := parser.NewParser(&l)
	   	env := object.NewEnv()
	   	e := evaluator.Eval(p.ParseProg(), env)
	   	fmt.Println(e)
	   	//fmt.Printf("AST:\n%v\n", p.ParseProg().ToString())

	   	if len(p.GetErrors()) > 0 {
	   		var errs string

	   		for _, err := range p.GetErrors() {
	   			errs += fmt.Sprintf("%s\n", err)
	   		}

	   		log.Warnln(errs)
	   	}
	*/

	//var name string = "পলাশ"
	//fmt.Println(len(name))
	//fmt.Println(string([]rune(name)[3]))

	//fmt.Println(name[12])

	args := os.Args[1:]

	if len(args) >= 1 {
		filename := args[0]
		_, err := os.Stat(filename)

		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("File `%s` does not exist!", filename)
		}

		f, err := os.ReadFile(filename)

		if err != nil {
			log.Fatalf("Cannot read `%s`", filename)
		}

		//fmt.Println(string(f))

		lx := lexer.NewLexer(string(f))
		ps := parser.NewParser(&lx)
		at := ps.ParseProg()

		if len(ps.GetErrors()) != 0 {
			repl.ShowParseErrors(os.Stdin, ps.GetErrors())
			log.Fatalf("fix above mentioned errors first!")
		}
		env := object.NewEnv()
		evd := evaluator.Eval(at, env)

		if evd != nil {
			fmt.Println(evd.Inspect())
		}

		//fmt.Println(args[0])

	}

	startRepl := true

	////
	//var name = "Palash Bauri"
	//var age = 20
	//fmt.Println(age + 2001)
	////

	if startRepl {
		user, err := user.Current()

		if err != nil {
			panic(err)
		}

		fmt.Printf("Hey, %s\n", user.Username)

		repl.Repl(os.Stdin, os.Stdout)
	}

}
