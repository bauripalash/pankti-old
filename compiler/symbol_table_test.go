package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	ex := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	g := NewSymbolTable()

	a := g.Define("a")
	if a != ex["a"] {
		t.Errorf("expected a=%+v, got=%+v", ex["a"], a)
	}
}

func TestResolveGlobal(t *testing.T) {
	g := NewSymbolTable()
	g.Define("a")

	ex := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
	}

	for _, s := range ex {
		res, ok := g.Resolve(s.Name)
		if !ok {
			t.Errorf("name %s not found", s.Name)
			continue
		}

		if res != s {
			t.Errorf("ex %s to resolve to %+v, got => %+v", s.Name, s, res)
		}

	}
}

func TestResolveLocal(t *testing.T) {
	g := NewSymbolTable()
	g.Define("a")
	g.Define("b")

	fl := NewEncolsedSymbolTable(g)
	fl.Define("c")
	fl.Define("d")

	tests := []struct {
		table  *SymbolTable
		expSym []Symbol
	}{
		{
			fl,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},

				{Name: "b", Scope: GlobalScope, Index: 1},

				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expSym {
			result, ok := tt.table.Resolve(sym.Name)

			if !ok {
				t.Errorf("name not found %s", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v; got %+v", sym.Name, sym, result)
			}
		}
	}
}
