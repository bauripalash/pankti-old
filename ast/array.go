package ast

import (
	"bytes"
	"go.cs.palashbauri.in/pankti/token"
	"strings"
)

// Arrays
type ArrLit struct {
	Token token.Token
	Elms  []Expr
}

func (ar *ArrLit) exprNode()        {}
func (ar *ArrLit) TokenLit() string { return ar.Token.Literal }
func (ar *ArrLit) String() string {
	var out bytes.Buffer

	es := []string{}

	for _, e := range ar.Elms {
		es = append(es, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")

	return out.String()
}

// Index Expression -> ARRAY[123]

type IndexExpr struct {
	Token token.Token
	Left  Expr
	Index Expr
}

func (ie *IndexExpr) exprNode()        {}
func (ie *IndexExpr) TokenLit() string { return ie.Token.Literal }
func (ie *IndexExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
