package evaluator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

func evalShowStmt(
	args []object.Obj,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	output := []string{}

	for _, item := range args {
		output = append(output, item.Inspect())
		//buff.Write([]byte(item.Inspect()))
	}

	if isGui {
		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			log.Fatalf(err.Error())
		}

		os.Stdout = w

		outC := make(chan string)

		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		fmt.Println(strings.Join(output, ""))
		w.Close()
		os.Stdout = oldStdout
		out := <-outC

		printBuff.Write([]byte(out))
	} else {
		fmt.Println(strings.Join(output, ""))
	}

	return NULL
}

func evalIncludeStmt(
	in *ast.IncludeStmt,
	e *object.Env,
	eh *ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) (*object.Env, object.Obj) {
	rawFilename := Eval(in.Filename, e, *eh, printBuff, isGui)
	enx := object.NewEnv()

	if rawFilename.Type() != object.STRING_OBJ {
		return enx, NewErr(
			rawFilename.GetToken(),
			eh,
			true,
			"include filename is invalid %s",
			rawFilename.Inspect(),
		)

	}

	includeFilename := rawFilename.(*object.String).Value

	_, err := os.Stat(includeFilename)

	if errors.Is(err, fs.ErrNotExist) {
		return enx, NewErr(
			in.Token,
			eh,
			true,
			"%s include file doesnot exists",
			includeFilename,
		)

	}

	fdata, err := os.ReadFile(includeFilename)

	if err != nil {
		return enx, NewErr(
			rawFilename.GetToken(),
			eh,
			true,
			"Failed to read include file %s",
			includeFilename,
		)

	}

	l := lexer.NewLexer(string(fdata))
	p := parser.NewParser(&l)
	ex := object.NewEnv()
	prog := p.ParseProg()
	Eval(prog, ex, *eh, printBuff, isGui)
	//fmt.Println(evd.Type())

	if len(p.GetErrors()) != 0 {
		for _, e := range p.GetErrors() {
			fmt.Println(e.String())
		}

		return enx, NewErr(
			rawFilename.GetToken(),
			eh,
			true,
			"Include file contains parsing errors",
		)
	}

	return ex, &object.Null{}

}

func evalBlockStmt(
	block *ast.BlockStmt,
	env *object.Env,
	eh *ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	var res object.Obj

	for _, stmt := range block.Stmts {
		res = Eval(stmt, env, *eh, printBuff, isGui)

		//fmt.Println("E_BS=> " , res)

		if res != nil {
			rtype := res.Type()
			if rtype == object.RETURN_VAL_OBJ || rtype == object.ERR_OBJ {
				//fmt.Println("RET => " ,  res)
				return res
			}
		}
	}
	//fmt.Println("EBS 2=>" ,res)
	return res
}
