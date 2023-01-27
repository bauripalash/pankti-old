package vm

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/compiler"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

func init() {
	log.SetLevel(log.ErrorLevel)
}

func parse(i string) *ast.Program {
	l := lexer.NewLexer(i)
	p := parser.NewParser(&l)
	return p.ParseProg()
}

func testNumberObject(t *testing.T, ex number.Number, got object.Obj) error {
	result, ok := got.(*object.Number)
	//	t.Log(result.Value.Value.String())
	//	t.Log(ex.Value.String())
	if !ok {
		return fmt.Errorf("Object is not a number. got=>%T , (%+v)", got, got)
	}

	if ex.Value.String() != result.Value.Value.String() {
		return fmt.Errorf("Object Wrong value got =>%v , Wanted => %v ",
			result.Value.Value, ex.Value,
		)
	}

	return nil
}

type vmTestCase struct {
	input string
	exp   interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		prog := parse(tt.input)
		comp := compiler.NewCompiler()
		err := comp.Compile(prog)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := NewVM(*comp.ByteCode())
		err = vm.Run()
		t.Logf("%v", vm.instructions)
		if err != nil {
			t.Fatalf("vm error : %s", err)
		}

		sElm := vm.StackTop()

		testExpectedObj(t, tt.exp, sElm)
	}
}

func testExpectedObj(t *testing.T,
	exp interface{},
	got object.Obj,
) {
	t.Helper()
	//t.Log(exp)
	//t.Log(got)
	switch exp := exp.(type) {

	case number.Number:
		err := testNumberObject(t, exp, got)
		if err != nil {
			t.Errorf("Number Object failed: %s", err)
		}
	}

}

func TestNumber(t *testing.T) {
	tests := []vmTestCase{
		{"1", number.MakeInt(1)},
		{"100", number.MakeInt(100)},
		{"1 + 2", number.MakeInt(3)},
	}

	//t.Log(tests)
	runVmTests(t, tests)
}
