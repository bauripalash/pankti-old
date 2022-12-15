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

func (s *StringLit) exprNode()        {}
func (s *StringLit) TokenLit() string { return s.Token.Literal }
func (s *StringLit) String() string   { return s.Token.Literal }

// =============================================================

// ================= Identifier ================================
type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) exprNode() {}
func (id *Identifier) TokenLit() string {
	return id.Token.Literal
}

func (id *Identifier) String() string {

	return id.Value
}

// =============================================================

// ================ Numbers ====================================
// Examples -> 100 , 200 , 23.3
type NumberLit struct {
	Token token.Token
	Value number.Number
	IsInt bool
}

func (nl *NumberLit) exprNode()        {}
func (nl *NumberLit) TokenLit() string { return nl.Token.Literal }
func (nl *NumberLit) String() string   { return nl.Token.Literal }

// ============================================================

// ================= Boolean ==================================
// Examples -> True , Sotto , False....
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) exprNode()        {}
func (b *Boolean) TokenLit() string { return b.Token.Literal }
func (b *Boolean) String() string   { return b.Token.Literal }

// ===========================================================
