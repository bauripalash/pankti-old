package parser

import (
	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) parseGroupedExpr() ast.Expr {

	p.nextToken()
	exp := p.parseExpr(LOWEST)

	if !p.peek(token.RPAREN) {
		return nil

	}

	return exp
}

func (p *Parser) parseExpr(prec int) ast.Expr {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		p.noPrefixFunctionErr(p.curTok)
		return nil
	}

	leftExpr := prefix()

	for !p.isPeekToken(token.SEMICOLON) && prec < p.peekPrec() {
		infix := p.infixParseFns[p.peekTok.Type]

		if infix == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infix(leftExpr)
	}

	//fmt.Println(leftExpr)
	return leftExpr

}

func (p *Parser) parsePrefixExpr() ast.Expr {
	exp := &ast.PrefixExpr{
		Token: p.curTok,
		Op:    p.curTok.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpr(PREFIX)

	log.Info("PREFIX => ", exp.Token, exp.Right)
	return exp
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {

	exp := &ast.InfixExpr{
		Token: p.curTok,
		Op:    p.curTok.Literal,
		Left:  left,
	}

	prec := p.curPrec()
	p.nextToken()
	exp.Right = p.parseExpr(prec)

	log.Info("INFIX => ", exp.Left, exp.Op, exp.Right)

	return exp
}

func (p *Parser) parseExprList(endtok token.TokenType) []ast.Expr {
	list := []ast.Expr{}

	if p.isPeekToken(endtok) {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, p.parseExpr(LOWEST))

	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpr(LOWEST))
	}

	if !p.peek(endtok) {
		return nil
	}

	return list
}
