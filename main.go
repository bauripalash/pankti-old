//go:generate goversioninfo -64 -icon=windows/res/icon.ico -manifest=windows/res/pankti.exe.manifest
package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/compiler"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/parser"
	"go.cs.palashbauri.in/pankti/vm"
)

/*
import (

	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"

	"go.cs.palashbauri.in/pankti/cmd"
	"go.cs.palashbauri.in/pankti/constants"

	log "github.com/sirupsen/logrus"

)

const DEBUG = true
const IMPORT_PATH_ENV = constants.IMPORT_PATH_ENV
const WINDOWS_IMPORT_PATH = "??" //TODO: Test on Windows
const LINUX_IMPORT_PATH = ".pankti/stdlib/"
*/
func init() {
	//log.SetLevel(log.DebugLevel)

	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	//log.SetOutput(os.Stdout)

	//setImportPathEnv()
}

/*
func setImportPathEnv() {
	custom_import_path := os.Getenv(IMPORT_PATH_ENV)

	if len(custom_import_path) < 1 && !DEBUG {

		if runtime.GOOS == "linux" {
			p, err := os.UserHomeDir()
			if err != nil {
				return
			}

			os.Setenv(IMPORT_PATH_ENV, path.Join(p, LINUX_IMPORT_PATH))

		} else if runtime.GOOS == "windows" {

			p, err := os.UserHomeDir()
			if err != nil {
				return
			}
			os.Setenv(IMPORT_PATH_ENV, path.Join(p, WINDOWS_IMPORT_PATH))

		}

	} else {
		curdir, _ := os.Getwd()
		os.Setenv(IMPORT_PATH_ENV, path.Join(curdir, "/stdlib/x/"))
	}
}

func main() {
	is_noide := false
	bi, noerr := debug.ReadBuildInfo()
	if !noerr {
		return
	}
	for _, item := range bi.Settings {
		if item.Key == "-tags" && strings.Contains(item.Value, "noide") {
			is_noide = true
			break
		}
	}

	cmd.Execute(is_noide)
}
*/

func helloparser(src string) ast.Program {
	l := lexer.NewLexer(src)
	p := parser.NewParser(&l)
	return *p.ParseProg()

	//return ast.Program{}
}

type t struct {
	Op       code.OpCode
	Operands []int
	Br       int
}

func main() {
	//helloparser("1 + 2 * 3 * 4 / 5 * 6 + 7 - 8")
	h := helloparser("{1:2}")
	/*tx := t{
		Op: code.OpConstant,
		Operands: []int{10},
		Br: 2,
	}
	*/
	//def , _ := code.Lookup(byte(tx.Op))
	//ins := code.Make(tx.Op , tx.Operands...)
	//or , n := code.ReadOperands(def , ins[1:])
	//fmt.Println(or , n)
	//fmt.Println(code.Instruction(ins))
	n := compiler.NewCompiler()
	n.Compile(&h)
	fmt.Println(n.ByteCode().Instructions.String())
	v := vm.NewVM(*n.ByteCode())
	v.Run()
	//fmt.Println(n.ByteCode().Constants[len(n.ByteCode().Constants)-1])
	fmt.Println(n)

}
