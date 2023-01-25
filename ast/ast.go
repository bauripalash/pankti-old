package ast

import (
	"bytes"
	"go.cs.palashbauri.in/pankti/token"
)

type Node interface {
	TokenLit() string
	String() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

// Program - main entry point
type Program struct {
	Source string //Main Source code
	Stmts  []Stmt //List of all statements
}

func (p *Program) TokenLit() string {

	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLit()
	} else {
		return ""
	}

}

func (p *Program) String() string {

	var out bytes.Buffer

	for _, stmt := range p.Stmts {
		out.WriteString(stmt.String())
	}

	return out.String()

}

type Comment struct {
	Token token.Token
	Value string
}

func (*Comment) stmtNode()          {}
func (s *Comment) TokenLit() string { return s.Token.Literal }
func (s *Comment) String() string   { return s.Token.Literal }
