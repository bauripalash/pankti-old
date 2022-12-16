package parser

import (
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) parseHashLit() ast.Expr {

	hash := &ast.HashLit{Token: p.curTok}
	hash.Pairs = make(map[ast.Expr]ast.Expr)

	for !p.isPeekToken(token.RBRACE) {
		p.nextToken()
		k := p.parseExpr(LOWEST)

		if !p.peek(token.COLON) {
			return nil
		}

		p.nextToken()

		val := p.parseExpr(LOWEST)

		hash.Pairs[k] = val

		if !p.isPeekToken(token.RBRACE) && !p.peek(token.COMMA) {
			return nil
		}
	}

	if !p.peek(token.RBRACE) {
		return nil
	}
	return hash

}

func (p *Parser) parseIndexExpr(l ast.Expr) ast.Expr {
	e := &ast.IndexExpr{Token: p.curTok, Left: l}

	p.nextToken()

	e.Index = p.parseExpr(LOWEST)

	if !p.peek(token.RS_BRACKET) {
		return nil
	}

	return e
}

func (p *Parser) parseArrLit() ast.Expr {
	arr := &ast.ArrLit{Token: p.curTok}

	arr.Elms = p.parseExprList(token.RS_BRACKET)

	return arr
}
