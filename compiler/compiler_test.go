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
		{
			input:   "1<2",
			exConst: []interface{}{1, 2},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGT),
				code.Make(code.OpPop),
			},
		},
		{
			input:   "jodi (sotto) tahole 10 nahole sesh 3333",
			exConst: []interface{}{10, 3333},
			exIns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 11),
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},

		{
			input:   "jodi (sotto) tahole 10 nahole 20 sesh 3333",
			exConst: []interface{}{10, 20, 3333},
			exIns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 13),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},

		{
			input: `
			dhori a = 1
			dhori b = 2
			`,
			exConst: []interface{}{1, 2},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},

		{

			input: `
            dhori a = 1
			a 
			`,
			exConst: []interface{}{1},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:   `"hello"`,
			exConst: []interface{}{"hello"},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		}, {
			input:   `"hell" + "o"`,
			exConst: []interface{}{"hell", "o"},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),

				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:   "[]",
			exConst: []interface{}{},
			exIns: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:   "[1,2,3]",
			exConst: []interface{}{1, 2, 3},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},

		{
			input:   "{1:2 , 3 : 4}",
			exConst: []interface{}{1, 2, 3, 4},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
		{
			input:   "{1:2+3 , 4 : 5*6}",
			exConst: []interface{}{1, 2, 3, 4, 5, 6},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
		{
			input:   "[1,2,3][1+1]",
			exConst: []interface{}{1, 2, 3, 1, 1},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),

				code.Make(code.OpConstant, 1),

				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),

				code.Make(code.OpConstant, 3),

				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		}, {
			input:   "{1:2}[2-1]",
			exConst: []interface{}{1, 2, 2, 1},
			exIns: []code.Instructions{
				code.Make(code.OpConstant, 0),

				code.Make(code.OpConstant, 1),

				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),

				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
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
			t.Fatalf("test Instructions failed : %s", err)
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
		return fmt.Errorf("Wrong Ins len;\n W=>%q\n G=>%q\n", ctd, got)
	}

	for i, ins := range ctd {
		if got[i] != ins {
			return fmt.Errorf("wrong instructions at %d; W=>%q\n G=>%q\n", i, ctd, got)
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
		case string:
			if err := testStringObj(con, got[i]); err != nil {
				return fmt.Errorf("constant %d - testStringObj failed %s", i, err)
			}

		}
	}

	return nil
}

func testStringObj(ex string, got object.Obj) error {
	r, ok := got.(*object.String)
	if !ok {
		return fmt.Errorf("Object not string got=%T (%+v)", got, got)
	}

	if r.Value != ex {
		return fmt.Errorf("Obj Wrong Obj got=%q want %q", r.Value, ex)
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
