package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store  map[string]Symbol
	numDef int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string) Symbol {
	sm := Symbol{Name: name, Index: s.numDef, Scope: GlobalScope}
	s.store[name] = sm
	s.numDef++
	return sm
}

func (s *SymbolTable) Resolve(n string) (Symbol, bool) {
	r, ok := s.store[n]
	return r, ok
}
