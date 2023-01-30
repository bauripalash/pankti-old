package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer  *SymbolTable
	store  map[string]Symbol
	numDef int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func NewEncolsedSymbolTable(outer *SymbolTable) *SymbolTable {
	st := NewSymbolTable()
	st.Outer = outer
	return st
}

func (s *SymbolTable) Define(name string) Symbol {
	sm := Symbol{Name: name, Index: s.numDef}

	if s.Outer == nil {
		sm.Scope = GlobalScope
	} else {
		sm.Scope = LocalScope
	}

	s.store[name] = sm
	s.numDef++
	return sm
}

func (s *SymbolTable) Resolve(n string) (Symbol, bool) {
	r, ok := s.store[n]
	if !ok && s.Outer != nil {
		obj, ok := s.Outer.Resolve(n)
		return obj, ok
	}
	return r, ok
}
