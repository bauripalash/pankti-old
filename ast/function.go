package ast

import (
	"bytes"
	"strings"

	"go.cs.palashbauri.in/pankti/token"
)

type FunctionLit struct {
	Token  token.Token // The 'fn' token
	Params []*Identifier
	Body   *BlockStmt
}

func (_ *FunctionLit) exprNode()         {}
func (fl *FunctionLit) TokenLit() string { return fl.Token.Literal }
func (fl *FunctionLit) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLit())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpr struct {
	Token token.Token // The '(' token
	Func  Expr
	// Identifier or FunctionLiteral
	Args []Expr
}

func (_ *CallExpr) exprNode()         {}
func (ce *CallExpr) TokenLit() string { return ce.Token.Literal }
func (ce *CallExpr) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Args {
		args = append(args, a.String())

	}
	out.WriteString(ce.Func.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
