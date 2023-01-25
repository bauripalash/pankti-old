package ast

import (
	"bytes"

	"go.cs.palashbauri.in/pankti/token"
)

type IfExpr struct {
	Token     token.Token
	Cond      Expr
	TrueBlock *BlockStmt
	ElseBlock *BlockStmt
}

func (*IfExpr) exprNode()          {}
func (i *IfExpr) TokenLit() string { return i.Token.Literal }
func (i *IfExpr) String() string {

	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Cond.String())
	out.WriteString(" ")
	out.WriteString(i.TrueBlock.String())
	if i.ElseBlock != nil {

		out.WriteString("else ")
		out.WriteString(i.ElseBlock.String())
	}
	return out.String()

}

type WhileExpr struct {
	Token     token.Token
	Cond      Expr
	StmtBlock *BlockStmt
}

func (*WhileExpr) exprNode()          {}
func (w *WhileExpr) TokenLit() string { return w.Token.Literal }
func (w *WhileExpr) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(w.Cond.String())
	out.WriteString(" ")
	out.WriteString(w.StmtBlock.String())

	return out.String()
}
