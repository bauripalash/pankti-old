package compiler

import (
	"fmt"
	"testing"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

type cTestCase struct {
	input   string
	exConst []interface{}
	exIns   []code.Instructions
}

func parse(i string) *ast.Program {
	l := lexer.NewLexer(i)
	p := parser.NewParser(&l)
	return p.ParseProg()
}

func TestIntArithmetic(t *testing.T) {
	tests := []cTestCase{
		{
			input:   "1+2",
			exConst: []interface{}{1, 2},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:   "sotto",
			exConst: []interface{}{},
			exIns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},

		{
			input:   "sotto",
			exConst: []interface{}{},
			exIns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},

		{
			input:   "mittha",
			exConst: []interface{}{},
			exIns: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
	}

	runCTests(t, tests)
}

func runCTests(t *testing.T, tests []cTestCase) {
	t.Helper()

	for _, tt := range tests {
		prog := parse(tt.input)
		compiler := NewCompiler()
		err := compiler.Compile(prog)

		if err != nil {
			t.Fatalf("compiler error : %s", err)
		}

		bc := compiler.ByteCode()

		err = testInstructions(tt.exIns, bc.Instructions)

		if err != nil {
			t.Fatalf("test Instructions faile : %s", err)
		}

		err = testConsts(t, tt.exConst, bc.Constants)

		if err != nil {
			t.Fatalf("testConsts failed : %s", err)
		}
	}
}

func testInstructions(ex []code.Instructions, got code.Instructions) error {
	ctd := concatIns(ex)

	if len(got) != len(ctd) {
		return fmt.Errorf("Wrong Ins len; W=>%q G=>%q", ctd, got)
	}

	for i, ins := range ctd {
		if got[i] != ins {
			return fmt.Errorf("wrong instructions at %d; W=>%q G=>%q", i, ctd, got)
		}
	}

	return nil
}

func testConsts(t *testing.T, ex []interface{}, got []object.Obj) error {

	if len(ex) != len(got) {
		return fmt.Errorf("wrong len of Constants; W=>%d , G=>%d", len(ex), len(got))
	}

	for i, con := range ex {
		switch con := con.(type) {
		case *number.Number:
			err := testIntObj(*con, got[i])

			if err != nil {
				return fmt.Errorf("const %d - testIntObj failed %s", i, err)
			}
		case *object.Boolean:
			if err := testBoolObj(*con, got[i]); err != nil {
				return fmt.Errorf("const %d testBoolObj failed %s", i, err)
			}

		}
	}

	return nil
}

func testIntObj(ex number.Number, got object.Obj) error {
	result, ok := got.(*object.Number)

	if !ok {
		return fmt.Errorf("not int %T (%+v)", got, got)
	}

	if result.Value.String() != ex.String() {
		return fmt.Errorf("obj wrong value W=>%s G=>%s", ex.Value.String(), result.Value.String())
	}

	return nil
}

func testBoolObj(ex object.Boolean, got object.Obj) error {

	result, ok := got.(*object.Boolean)

	if !ok {
		return fmt.Errorf("not int %T (%+v)", got, got)
	}

	if result.Value != ex.Value {
		return fmt.Errorf("wrong bool value; W=>%T , G=>%T", ex.Value, result.Value)
	}

	return nil

}

func concatIns(i []code.Instructions) code.Instructions {
	r := code.Instructions{}

	for _, ins := range i {
		r = append(r, ins...)
	}

	return r
}
