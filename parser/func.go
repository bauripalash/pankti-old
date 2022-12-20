package parser

import (
	//"fmt"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) parseFunc() ast.Expr {

	if !p.peek(token.FUNC) {

		return nil

	}

	fl := &ast.FunctionLit{Token: p.curTok}
	//fmt.Println(fl.Token)
	if !p.peek(token.LPAREN) {
		return nil
	}

	fl.Params = p.parseFuncParams()

	//if !p.peek(token.LBRACE) {
	//	return nil
	//}

	fl.Body = p.parseBlockStmt(token.END)

	log.Info("FN EXPR => ", fl.Body.String())

	return fl
}

func (p *Parser) parseFuncParams() []*ast.Identifier {
	ids := []*ast.Identifier{}

	if p.isPeekToken(token.RPAREN) {
		p.nextToken()
		return ids
	}

	p.nextToken()

	id := &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	ids = append(ids, id)

	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		id := &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
		ids = append(ids, id)
	}

	if !p.peek(token.RPAREN) {
		return nil
	}

	log.Info("FUNC PARAMS => ", ids)
	return ids
}

func (p *Parser) parseCallExpr(function ast.Expr) ast.Expr {
	exp := &ast.CallExpr{Token: p.curTok, Func: function}
	exp.Args = p.parseExprList(token.RPAREN)
	return exp
}
