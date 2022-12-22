package ast

import (
	"bytes"

	"go.cs.palashbauri.in/pankti/token"
)

// Prefix Expression

type PrefixExpr struct {
	Token token.Token
	Op    string
	Right Expr
}

func (_ *PrefixExpr) exprNode()           {}
func (pref *PrefixExpr) TokenLit() string { return pref.Token.Literal }
func (pref *PrefixExpr) String() string {

	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pref.Op)
	out.WriteString(pref.Right.String())
	out.WriteString(")")
	return out.String()

}

//Infix Expression

type InfixExpr struct {
	Token token.Token
	Left  Expr
	Op    string
	Right Expr
}

func (_ *InfixExpr) exprNode()          {}
func (inf *InfixExpr) TokenLit() string { return inf.Token.Literal }
func (inf *InfixExpr) String() string {

	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(inf.Left.String())
	out.WriteString(" " + inf.Op + " ")
	out.WriteString(inf.Right.String())
	out.WriteString(")")
	return out.String()

}
