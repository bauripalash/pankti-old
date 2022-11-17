package ast

import (
	"bytes"
	"strings"

	"bauri.palash/pankti/number"
	"bauri.palash/pankti/token"
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

//Program - main entry point

type Program struct {
	Source string
	Stmts  []Stmt
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

func (s *Comment) stmtNode()        {}
func (s *Comment) TokenLit() string { return s.Token.Literal }
func (s *Comment) String() string   { return s.Token.Literal }

type StringLit struct {
	Token token.Token
	Value string
}

func (s *StringLit) exprNode()        {}
func (s *StringLit) TokenLit() string { return s.Token.Literal }
func (s *StringLit) String() string   { return s.Token.Literal }

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

//Hash

type HashLit struct {
	Token token.Token
	Pairs map[Expr]Expr
}

func (hl *HashLit) exprNode()        {}
func (hl *HashLit) TokenLit() string { return hl.Token.Literal }
func (hl *HashLit) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// let statement

type LetStmt struct {
	Token token.Token
	Name  Identifier
	Value Expr
}

func (lst *LetStmt) stmtNode() {}
func (lst *LetStmt) TokenLit() string {
	return lst.Token.Literal
}

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

//Return statement

type ReturnStmt struct {
	Token     token.Token
	ReturnVal Expr
}

func (r *ReturnStmt) stmtNode() {}
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

type ShowStmt struct{
    Token token.Token
    Value []Expr
}

func (s *ShowStmt) stmtNode() {}
func (s *ShowStmt) TokenLit() string {
    return s.Token.Literal
}

func (s *ShowStmt) String() string{
    var out bytes.Buffer

    out.WriteString(s.TokenLit() + " ")
    if s.Value != nil{
        for _,itme := range s.Value{
            out.WriteString(itme.String())
            out.WriteString(",")
        }
    }

    out.WriteString(";")
    return out.String()

}

// Expression Statement
type ExprStmt struct {
	Token token.Token
	Expr  Expr
}

func (e *ExprStmt) stmtNode() {}
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

/*
func (e *ExprStmt) String() string {
	//fmt.Println(e.Expr.TokenLit())
	if e.Expr != nil {
		return e.Expr.ToString()
	} else {

		return ""
	}
}
*/

// Identifier Expression

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

type NumberLit struct {
	Token token.Token
	Value number.Number
	IsInt bool
}

func (nl *NumberLit) exprNode()        {}
func (nl *NumberLit) TokenLit() string { return nl.Token.Literal }
func (nl *NumberLit) String() string   { return nl.Token.Literal }

// Prefix Expression

type PrefixExpr struct {
	Token token.Token
	Op    string
	Right Expr
}

func (pref *PrefixExpr) exprNode()        {}
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

func (inf *InfixExpr) exprNode()        {}
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

// Boolean Expression
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) exprNode()        {}
func (b *Boolean) TokenLit() string { return b.Token.Literal }
func (b *Boolean) String() string   { return b.Token.Literal }

type IncludeStmt struct {
	Token    token.Token
	Filename Expr
}

func (is *IncludeStmt) stmtNode()        {}
func (is *IncludeStmt) TokenLit() string { return is.Token.Literal }
func (is *IncludeStmt) String() string   { return is.Token.Literal }

type BlockStmt struct {
	Token token.Token
	Stmts []Stmt
}

func (bs *BlockStmt) stmtNode()        {}
func (bs *BlockStmt) TokenLit() string { return bs.Token.Literal }
func (bs *BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range bs.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}

type IfExpr struct {
	Token     token.Token
	Cond      Expr
	TrueBlock *BlockStmt
	ElseBlock *BlockStmt
}

func (i *IfExpr) exprNode()        {}
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

func (w *WhileExpr) exprNode()        {}
func (w *WhileExpr) TokenLit() string { return w.Token.Literal }
func (w *WhileExpr) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(w.Cond.String())
	out.WriteString(" ")
	out.WriteString(w.StmtBlock.String())

	return out.String()
}

type FunctionLit struct {
	Token  token.Token // The 'fn' token
	Params []*Identifier
	Body   *BlockStmt
}

func (fl *FunctionLit) exprNode()        {}
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

func (ce *CallExpr) exprNode()        {}
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


