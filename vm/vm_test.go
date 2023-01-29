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

func testBoolObj(ex bool, got object.Obj) error {
	//t.Helper()

	r, ok := got.(*object.Boolean)

	if !ok {
		return fmt.Errorf("not bool got=>%T (%+v)", got, got)
	}

	if r.Value != ex {
		return fmt.Errorf("Bool obj Value mismatch W=>%v , G=>%v", r.Value, ex)
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
		if err != nil {
			t.Fatalf("vm error : %s", err)
		}

		sElm := vm.LastPoppedStackItem()

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
	case bool:
		t.Log(exp)
		t.Log(got)
		if err := testBoolObj(exp, got); err != nil {
			t.Errorf("Bool obj failed : %s", err)
		}
	case *object.Null:
		if got != Null {
			t.Errorf("Obj not null : %T (%+v)", got, got)
		}
	case string:
		if err := testStringObj(exp, got); err != nil {
			t.Errorf("String Obj Failed : %s", err)
		}
	case []number.Number:
		arr, ok := got.(*object.Array)

		if !ok {
			t.Errorf("OBJ not array : %T (%+v)", got, got)
			return
		}

		if len(arr.Elms) != len(exp) {
			t.Errorf("wrong number of elems W=%d G=%d", len(exp), len(arr.Elms))
			return
		}

		for i, ee := range exp {
			if err := testNumberObject(t, ee, arr.Elms[i]); err != nil {
				t.Errorf("testIntObj failed %s", err)
			}
		}
	case map[object.HashKey]number.Number:
		hash, ok := got.(*object.Hash)
		//ok := true
		//hash := object.Hash{ Pairs : make(map[object.HashKey]object.HashPair)}
		t.Log("---")
		t.Log(hash.Pairs)
		t.Log("---")
		if !ok {
			t.Errorf("object is not hash %T", got)
		}

		if len(hash.Pairs) != len(exp) {
			t.Errorf("hash has wrong number of pairs W=>%d G=>%d", len(exp), len(hash.Pairs))
		}

		for expKey, expValue := range exp {
			t.Log(expKey)
			pair, ok := hash.Pairs[expKey]
			if !ok {
				t.Errorf("No Pair for given key")
			}
			t.Logf("%v -> %v", expValue, pair)

		}

	}

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

func TestNumber(t *testing.T) {
	tests := []vmTestCase{
		/*		{"1", number.MakeInt(1)},
				{"100", number.MakeInt(100)},
				{"1 + 2", number.MakeInt(3)},
				{"1-3", number.MakeInt(-2)},

				{"2*3", number.MakeInt(6)},

				{"6/2", number.MakeInt(3)},

				{"sotto", true},
				{"mittha", false},
				{"1 == 2", false},
				{"jodi (sotto) tahole 10 nahole sesh", 10},
				{"jodi (sotto) tahole 10 nahole 20 sesh;", 10},
				{"jodi (mittha) tahole 10 nahole 20 sesh ", 20},

				{"jodi (1<2) tahole 10 nahole 20 sesh ", 10},
				{"jodi ( 1 > 2) tahole 10 nahole sesh", Null},
				{"!(jodi (mittha) tahole 5 nahole sesh )", true},
				{`dhori a = 1
				a`, 1},
				{`"hell" + "o"`, "hello"},
				{"[]", []int{}},
				{"[1 , 2 ,3]", []int{1, 2, 3}},
		*/ //	{"{}", map[object.HashKey]number.Number{}},
		{`{"1":2}`, map[object.HashKey]number.Number{
			(&object.String{Value: "1"}).HashKey(): number.MakeFloat(2),
		}},
		{"[1,2,3][1]", 2},
		{"[[1,1,1]][0][0]", 1},
		{"[][0]", Null},
		{`{"a" : 1 , "b" : 2}["a"]`, 1},
	}

	//t.Log(tests)
	runVmTests(t, tests)
}
