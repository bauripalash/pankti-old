package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"bauri.palash/pankti/ast"
	"bauri.palash/pankti/number"
	"bauri.palash/pankti/token"
	"strings"
)

const (
	INT_OBJ        = "INTEGER"
	FLOAT_OBJ      = "FLOAT"
	BOOL_OBJ       = "BOOLEAN"
	RETURN_VAL_OBJ = "RETURN_VAL"
	NULL_OBJ       = "NIL"
	ERR_OBJ        = "ERROR"
	FUNC_OBJ       = "FUNCTION"
	STRING_OBJ     = "STRING"
	BUILTIN_OBJ    = "BUILTIN"
	ARRAY_OBJ      = "ARRAY"
	HASH_OBJ       = "HASH"
	NUM_OBJ        = "NUM"
	INCLUDE_OBJ    = "INCLUDE"
    SHOW_OBJ       = "SHOW"
)

type BuiltInFunc func(args ...Obj) Obj

type ObjType string

type Obj interface {
	Type() ObjType
	Inspect() string
	GetToken() token.Token
}

type Builtin struct {
	Fn    BuiltInFunc
	Token token.Token
}

func (b *Builtin) Type() ObjType         { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string       { return "builtin function" }
func (b *Builtin) GetToken() token.Token { return b.Token }

// Include

type IncludeObj struct {
	Filename string
	Token    token.Token
}

func (ib *IncludeObj) Type() ObjType         { return INCLUDE_OBJ }
func (ib *IncludeObj) Inspect() string       { return fmt.Sprintf("include %s", ib.Filename) }
func (ib *IncludeObj) GetToken() token.Token { return ib.Token }

// Print Obj

type ShowObj struct{
    Value []string 
    Token token.Token
}

func (so *ShowObj) Type() ObjType { return  SHOW_OBJ }
func (so *ShowObj) Inspect() string { return strings.Join(so.Value, "\n") }
func (so *ShowObj ) GetToken() token.Token {return so.Token  }
func (so *ShowObj ) Print(b bytes.Buffer) {
   b.WriteString(so.Inspect()) 
   b.WriteString("\n")
}


//Arrays

type Array struct {
	Elms  []Obj
	Token token.Token
}

func (a *Array) Type() ObjType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	es := []string{}
	for _, e := range a.Elms {
		es = append(es, e.Inspect())
	}
	out.WriteString("[")

	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")
	return out.String()
}

func (a *Array) GetToken() token.Token {
	return a.Token
}

//Hash

type HashKey struct {
	Type  ObjType
	Value uint64
	Token token.Token
}

func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 2
	}

	return HashKey{Type: b.Type(), Value: val}
}

func (s *String) HashKey() HashKey {

	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (n *Number) HashKey() HashKey {
	h := fnv.New64a()
	if n.Value.IsInt {
		k := n.Value.Value.(*number.IntNumber).Value
		h.Write(k.Bytes())
		return HashKey{Type: n.Type(), Value: h.Sum64()}
	} else {
		k := n.Value.Value.(*number.FloatNumber).Value
		temp, _ := k.Int64()
		h.Write([]byte(k.String() + fmt.Sprintf("%d", temp)))
		return HashKey{Type: n.Type(), Value: h.Sum64()}
		//        h.Write()
	}

}

//Hash Pair { a : b}

type HashPair struct {
	Key   Obj
	Value Obj
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Token token.Token
}

func (h *Hash) Type() ObjType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, p := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s : %s", p.Key.Inspect(), p.Value.Inspect()))

	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (h *Hash) GetToken() token.Token {
	return h.Token
}

type Hashable interface {
	HashKey() HashKey
}

// Strings "I am a string"
type String struct {
	Value string
	Token token.Token
}

func (s *String) Type() ObjType         { return STRING_OBJ }
func (s *String) Inspect() string       { return s.Value }
func (s *String) GetToken() token.Token { return s.Token }

type Number struct {
	Value number.Number
	IsInt bool
	Token token.Token
}

func (num *Number) Type() ObjType {
	return NUM_OBJ
}

func (num *Number) Inspect() string {
	return fmt.Sprintf("%s", num.Value.Value.String())
}

func (num *Number) GetToken() token.Token {
	return num.Token
}

// Booleans true,false
type Boolean struct {
	Value bool
	Token token.Token
}

func (b *Boolean) Type() ObjType         { return BOOL_OBJ }
func (b *Boolean) Inspect() string       { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) GetToken() token.Token { return b.Token }

// NULL_OBJ
type Null struct{}

func (n *Null) Type() ObjType         { return NULL_OBJ }
func (n *Null) Inspect() string       { return "null" }
func (n *Null) GetToken() token.Token { return token.Token{} }

type ReturnValue struct {
	Value Obj
	Token token.Token
}

func (r *ReturnValue) Type() ObjType         { return RETURN_VAL_OBJ }
func (r *ReturnValue) Inspect() string       { return r.Value.Inspect() }
func (r *ReturnValue) GetToken() token.Token { return r.Token }

type Error struct {
	Msg string
}

func (e *Error) Type() ObjType         { return ERR_OBJ }
func (e *Error) Inspect() string       { return "ERR : " + e.Msg }
func (e *Error) GetToken() token.Token { return token.Token{} }

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStmt
	Env    *Env
	Token  token.Token
}

func (f *Function) Type() ObjType { return FUNC_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (f *Function) GetToken() token.Token {

	return f.Token
}
