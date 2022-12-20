package object

import (
	"bytes"
	"fmt"
	"strings"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/token"
)

type Function struct {
	Name   string
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

// Include

type IncludeObj struct {
	Filename string
	Token    token.Token
}

func (ib *IncludeObj) Type() ObjType { return INCLUDE_OBJ }

func (ib *IncludeObj) Inspect() string {

	return fmt.Sprintf("include %s", ib.Filename)
}

func (ib *IncludeObj) GetToken() token.Token { return ib.Token }

// Print Obj

type ShowObj struct {
	Value []string
	Token token.Token
}

func (so *ShowObj) Type() ObjType { return SHOW_OBJ }

func (so *ShowObj) Inspect() string {
	return strings.Join(so.Value, "\n")
}

func (so *ShowObj) GetToken() token.Token { return so.Token }
func (so *ShowObj) Print(b bytes.Buffer) {
	b.WriteString(so.Inspect())
	b.WriteString("\n")
}

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
