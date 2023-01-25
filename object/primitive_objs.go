package object

import (
	"bytes"
	"fmt"
	"strings"

	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/token"
)

// Strings "I am a string"
type String struct {
	Value string
	Token token.Token
}

func (*String) Type() ObjType           { return STRING_OBJ }
func (s *String) Inspect() string       { return s.Value }
func (s *String) GetToken() token.Token { return s.Token }

type Number struct {
	Value number.Number
	IsInt bool
	Token token.Token
}

func (*Number) Type() ObjType {
	return NUM_OBJ
}

func (num *Number) Inspect() string {
	return num.Value.Value.String()
}

func (num *Number) GetToken() token.Token {
	return num.Token
}

// Booleans true,false
type Boolean struct {
	Value bool
	Token token.Token
}

func (*Boolean) Type() ObjType { return BOOL_OBJ }

func (b *Boolean) Inspect() string       { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) GetToken() token.Token { return b.Token }

type Break struct {
	Token     token.Token
	PrevValue Obj
}

func (*Break) Type() ObjType           { return BREAK_OBJ }
func (*Break) Inspect() string         { return "break" }
func (b *Break) GetToken() token.Token { return b.Token }
func (b *Break) Prev() Obj             { return b.PrevValue }

// NULL_OBJ
type Null struct{}

func (*Null) Type() ObjType         { return NULL_OBJ }
func (*Null) Inspect() string       { return "null" }
func (*Null) GetToken() token.Token { return token.Token{} }

//Arrays

type Array struct {
	Elms  []Obj
	Token token.Token
}

func (*Array) Type() ObjType { return ARRAY_OBJ }
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
