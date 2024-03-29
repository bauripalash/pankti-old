package ast

import (
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/token"
)

// ==================== Strings ===============================
// example -> "hello world"
type StringLit struct {
	Token token.Token
	Value string
}

func (*StringLit) exprNode()          {}
func (s *StringLit) TokenLit() string { return s.Token.Literal }
func (s *StringLit) String() string   { return s.Token.Literal }

// =============================================================

// ================= Identifier ================================
type Identifier struct {
	Token token.Token
	IsMod bool
	Value string
}

func (*Identifier) exprNode() {}
func (id *Identifier) TokenLit() string {
	return id.Token.Literal
}

func (id *Identifier) String() string {

	return id.Value
}

type IncludeId struct {
	Token token.Token
	Value string
}

func (*IncludeId) exprNode()           {}
func (ii *IncludeId) TokenLit() string { return ii.Token.Literal }
func (ii *IncludeId) String() string {
	return ii.Value
}

type IncludeExpr struct {
	Token    token.Token
	Filename Expr
}

func (*IncludeExpr) exprNode()           {}
func (ix *IncludeExpr) TokenLit() string { return ix.Token.Literal }
func (ix *IncludeExpr) String() string   { return ix.Token.Literal }

// =============================================================

// ================ Numbers ====================================
// Examples -> 100 , 200 , 23.3
type NumberLit struct {
	Token token.Token
	Value number.Number
	IsInt bool
}

func (*NumberLit) exprNode()           {}
func (nl *NumberLit) TokenLit() string { return nl.Token.Literal }
func (nl *NumberLit) String() string   { return nl.Token.Literal }

// ============================================================

// ================= Boolean ==================================
// Examples -> True , Sotto , False....
type Boolean struct {
	Token token.Token
	Value bool
}

func (*Boolean) exprNode()          {}
func (b *Boolean) TokenLit() string { return b.Token.Literal }
func (b *Boolean) String() string   { return b.Token.Literal }

// ===========================================================

type Break struct {
	Token token.Token
	Value string
}

func (*Break) exprNode()          {}
func (b *Break) TokenLit() string { return b.Token.Literal }
func (b *Break) String() string   { return b.Token.Literal }
