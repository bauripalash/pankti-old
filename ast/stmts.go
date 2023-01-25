package ast

import (
	"bytes"

	"go.cs.palashbauri.in/pankti/token"
)

// ====================== Let / Dhori Statment ========================
// Example -> dhori age = 20
type LetStmt struct {
	Token token.Token
	Name  Identifier
	Value Expr
}

func (*LetStmt) stmtNode()            {}
func (lst *LetStmt) TokenLit() string { return lst.Token.Literal }

func (lst *LetStmt) String() string {

	var out bytes.Buffer

	out.WriteString(lst.TokenLit() + " ")
	out.WriteString(lst.Name.String())
	out.WriteString(" = ")

	if lst.Value != nil {
		out.WriteString(lst.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ==================================================================

// ==================== Return / Ferau statement ============
// Example -> return true
type ReturnStmt struct {
	Token     token.Token
	ReturnVal Expr
}

func (*ReturnStmt) stmtNode() {}
func (r *ReturnStmt) TokenLit() string {
	return r.Token.Literal
}

func (r *ReturnStmt) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLit() + " ")
	if r.ReturnVal != nil {
		out.WriteString(r.ReturnVal.String())
	}
	out.WriteString(";")
	return out.String()
}

// ========================================================

// ================== Print / Show / Dekhau Statment ============
// Example -> show("Hello World")
type ShowStmt struct {
	Token token.Token
	Value []Expr
}

func (*ShowStmt) stmtNode() {}
func (s *ShowStmt) TokenLit() string {
	return s.Token.Literal
}

func (s *ShowStmt) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLit() + " ")
	if s.Value != nil {
		for _, itme := range s.Value {
			out.WriteString(itme.String())
			out.WriteString(",")
		}
	}

	out.WriteString(";")
	return out.String()

}

// ===============================================================

type IncludeStmt struct {
	Token    token.Token
	Filename Expr
}

func (*IncludeStmt) stmtNode()           {}
func (is *IncludeStmt) TokenLit() string { return is.Token.Literal }
func (is *IncludeStmt) String() string   { return is.Token.Literal }

type BlockStmt struct {
	Token token.Token
	Stmts []Stmt
}

func (*BlockStmt) stmtNode()           {}
func (bs *BlockStmt) TokenLit() string { return bs.Token.Literal }
func (bs *BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range bs.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}

// Expression Statement
type ExprStmt struct {
	Token token.Token
	Expr  Expr
}

func (*ExprStmt) stmtNode() {}
func (e *ExprStmt) TokenLit() string {
	return e.Token.Literal
}

func (e *ExprStmt) String() string {
	//fmt.Println(e.Expr.TokenLit())
	if e.Expr != nil {
		return e.Expr.String()
	} else {

		return ""
	}
}
