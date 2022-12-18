package parser

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) parseShowStmt() *ast.ShowStmt {

	stmt := &ast.ShowStmt{Token: p.curTok}
	p.nextToken()
	stmt.Value = p.parseExprList(token.RPAREN)
	//p.nextToken()
	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	log.Info(fmt.Sprintf("SHOW STMT => %v\n", stmt))

	return stmt
}

func (p *Parser) parseIncludeStmt() *ast.IncludeStmt {
	stmt := &ast.IncludeStmt{Token: p.curTok}
	p.nextToken()

	stmt.Filename = p.parseExpr(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	log.Info(
		fmt.Sprintf(
			"INCLUDE => FNAME=>%s || FNAME_TYPE=>%s",
			stmt.Filename,
			stmt,
		),
	)

	return stmt
}

func (p *Parser) parseIncludeExpr() ast.Expr {
	exp := &ast.IncludeExpr{Token: p.curTok}
	p.nextToken()

	exp.Filename = p.parseExpr(LOWEST)
	log.Info(fmt.Sprintf(
		"INCLUDE EXPR => FNAME->%s",
		exp.Filename,
	))

	return exp
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	//LET <IDENTIFIER> <EQUAL_SIGN> <EXPRESSION>
	stmt := &ast.LetStmt{Token: p.curTok}

	if !p.peek(token.IDENT) {
		return nil
	}
	isModId := len(strings.Split(p.curTok.Literal, ".")) == 2

	stmt.Name = ast.Identifier{
		Token: p.curTok,
		IsMod: isModId,
		Value: p.curTok.Literal,
	}

	if !p.peek(token.EQ) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpr(LOWEST)

	for p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	log.Info(fmt.Sprintf("LET STMT => %v\n", stmt))
	return stmt

}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curTok}

	p.nextToken()

	stmt.ReturnVal = p.parseExpr(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	log.Info(fmt.Sprintf("RETURN STMT => %v\n", stmt))

	return stmt

}

func (p *Parser) parseBlockStmt(eT token.TokenType) *ast.BlockStmt {
	bs := &ast.BlockStmt{Token: p.curTok, Stmts: []ast.Stmt{}}

	//	bs.Stmts = []ast.Stmt{}

	p.nextToken()

	for !p.isCurToken(eT) && !p.isCurToken(token.EOF) {
		s := p.parseStmt()
		if s != nil {
			bs.Stmts = append(bs.Stmts, s)
		}
		p.nextToken()
	}
	//fmt.Println("BS=> " , bs)

	return bs
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	//fmt.Println(p.curTok)
	stmt := &ast.ExprStmt{Token: p.curTok}

	stmt.Expr = p.parseExpr(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}
	//fmt.Println("expr stmt->>>" , stmt)
	return stmt
}
