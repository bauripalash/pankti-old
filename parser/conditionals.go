package parser

import (
	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) parseIfExpr() ast.Expr {

	exp := &ast.IfExpr{Token: p.curTok}
	has_else := false
	if !p.peek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Cond = p.parseExpr(LOWEST)

	if !p.peek(token.RPAREN) {
		return nil
	}
	// jodi (sotto) tahole { "hello" }
	if !p.peek(token.TAHOLE) {
		return nil
	}
	p.nextToken()
	tb := &ast.BlockStmt{Token: p.curTok, Stmts: []ast.Stmt{}}
	eb := &ast.BlockStmt{Token: p.curTok, Stmts: []ast.Stmt{}}

	for !p.isCurToken(token.ELSE) && !p.isCurToken(token.EOF) {
		s := p.parseStmt()
		//if  {
		tb.Stmts = append(tb.Stmts, s)
		//}
		p.nextToken()
	}

	p.nextToken()

	if !p.isCurToken(token.END) && !p.isCurToken(token.EOF) {
		s := p.parseStmt()
		//if s != nil {
		eb.Stmts = append(eb.Stmts, s)
		//}
		p.nextToken()
	}

	exp.TrueBlock = tb
	exp.ElseBlock = eb

	//p.nextToken()
	//exp.TrueBlock = p.parseBlockStmt(token.ELSE)
	//p.nextToken()
	//exp.ElseBlock = p.parseBlockStmt(token.END)

	if has_else {
		log.Info(
			"IF ELSE Expr => ",
			exp.Cond,
			exp.TrueBlock.String(),
			exp.ElseBlock.String(),
		)
	} else {
		log.Info("IF Expr => ", exp.Cond, exp.TrueBlock.String())
	}

	return exp
}

func (p *Parser) parseWhileExpr() ast.Expr {

	exp := &ast.WhileExpr{Token: p.curTok}

	if !p.peek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Cond = p.parseExpr(LOWEST)

	if !p.peek(token.RPAREN) {
		return nil
	}

	//if !p.peek(token.LBRACE) {
	//	return nil
	//}

	exp.StmtBlock = p.parseBlockStmt(token.END)

	return exp

}
