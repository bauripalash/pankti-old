package main_test

import (
	"bytes"
	"testing"

	"go.cs.palashbauri.in/pankti/evaluator"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

var src = `
ধরি ফিব = একটি কাজ(ক)

    যদি (ক < 0) তাহলে
        দেখাও("ভুল ইনপুট")
    নাহলে
        যদি(ক == 0) তাহলে
            ফেরাও(0)
        নাহলে
            যদি (ক == ১) তাহলে
                ফেরাও(১)
            নাহলে
                যদি (ক == ২) তাহলে
                    ফেরাও(১)
                নাহলে
                    ফেরাও (ফিব(ক-১) + ফিব(ক-২))
                শেষ
            শেষ
        শেষ
    শেষ
শেষ

ফিব(30)

`
var expected = "832040"

func TestX(t *testing.T) {

	l := lexer.NewLexer(src)
	p := parser.NewParser(&l)
	prog := p.ParseProg()
	eh := object.ErrorHelper{Source: src}
	printBuf := new(bytes.Buffer)
	env := object.NewEnvMap()
	ev := evaluator.Eval(prog, env, eh, printBuf, true)

	got := ev.Inspect()

	if got != expected {
		t.Errorf("Got %q; Wanted %q", got, expected)
	}

}

func BenchmarkX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := lexer.NewLexer(src)
		p := parser.NewParser(&l)
		prog := p.ParseProg()
		eh := object.ErrorHelper{Source: src}
		printBuf := new(bytes.Buffer)
		env := object.NewEnvMap()
		evaluator.Eval(prog, env, eh, printBuf, true)
	}
}
