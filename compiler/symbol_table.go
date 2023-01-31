package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
	FreeScope 	SymbolScope = "FREE"
	FuncScope 	SymbolScope	= "FUNCTION"
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
	FreeSymbols []Symbol
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	fs := []Symbol{}
	return &SymbolTable{store: s , FreeSymbols: fs}
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

func (s *SymbolTable) DefineFuncName(name string) Symbol{
	sm := Symbol{ Name: name , Index: 0 , Scope: FuncScope }
	s.store[name] = sm
	return sm
}

func (s *SymbolTable) defineFree(o Symbol) Symbol{
	s.FreeSymbols = append(s.FreeSymbols, o)
	sm := Symbol { Name: o.Name , Index: len(s.FreeSymbols) - 1 }
	sm.Scope = FreeScope
	s.store[o.Name] = sm
	return sm
	
}

func (s *SymbolTable) Resolve(n string) (Symbol, bool) {
	r, ok := s.store[n]
	if !ok && s.Outer != nil {
		r, ok = s.Outer.Resolve(n)
		if !ok{
			return r , ok
		}
		if r.Scope == GlobalScope{
			return r , ok
		}

		free := s.defineFree(r)
		return free , true
	}


	return r, ok
}
