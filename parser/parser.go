package parser

import (
	"fmt"
	"math/big"
	"pankti/ast"
	"pankti/errs"
	"pankti/lexer"
	"pankti/number"
	"pankti/token"

	log "github.com/sirupsen/logrus"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LTGT
	SUM
	PROD
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{

	token.EQEQ:       EQUALS,
	token.NOT_EQ:     EQUALS,
	token.LT:         LTGT,
	token.GT:         LTGT,
	token.GTE:        LTGT,
	token.LTE:        LTGT,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.DIV:        PROD,
	token.MUL:        PROD,
	token.LPAREN:     CALL,
	token.LS_BRACKET: INDEX,
}

type Parser struct {
	lx      *lexer.Lexer
	curTok  token.Token
	peekTok token.Token

	errs []errs.ParserError

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

func NewParser(l *lexer.Lexer) *Parser {

	p := &Parser{lx: l,
		errs: []errs.ParserError{},
	}

	//register prefix functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.regPrefix(token.IDENT, p.parseIdent)
	//	p.regPrefix(token.INT, p.parseIntegerLit)
	//	p.regPrefix(token.FLOAT, p.parseFloatLit)
	p.regPrefix(token.NUM, p.parseNumLit)
	p.regPrefix(token.MINUS, p.parsePrefixExpr)
	p.regPrefix(token.EXC, p.parsePrefixExpr)
	p.regPrefix(token.TRUE, p.parseBool)
	p.regPrefix(token.FALSE, p.parseBool)
	p.regPrefix(token.LPAREN, p.parseGroupedExpr)
	p.regPrefix(token.IF, p.parseIfExpr)
	p.regPrefix(token.WHILE, p.parseWhileExpr)
	p.regPrefix(token.EKTI, p.parseFunc)
	p.regPrefix(token.STRING, p.parseStringLit)
	p.regPrefix(token.LS_BRACKET, p.parseArrLit)
	p.regPrefix(token.LBRACE, p.parseHashLit)

	//register infix functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.regInfix(token.PLUS, p.parseInfixExpr)
	p.regInfix(token.MINUS, p.parseInfixExpr)
	p.regInfix(token.DIV, p.parseInfixExpr)
	p.regInfix(token.MUL, p.parseInfixExpr)
	p.regInfix(token.EQEQ, p.parseInfixExpr)
	p.regInfix(token.NOT_EQ, p.parseInfixExpr)
	p.regInfix(token.LT, p.parseInfixExpr)
	p.regInfix(token.GTE, p.parseInfixExpr)
	p.regInfix(token.GT, p.parseInfixExpr)
	p.regInfix(token.LTE, p.parseInfixExpr)
	p.regInfix(token.LPAREN, p.parseCallExpr)
	p.regInfix(token.LS_BRACKET, p.parseIndexExpr)

	p.nextToken()
	p.nextToken()
	//fmt.Println(p.curTok , p.peekTok)
	return p

}

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

func (p *Parser) parseStringLit() ast.Expr {
	//fmt.Println(p.curTok)
	return &ast.StringLit{Token: p.curTok, Value: p.curTok.Literal}
}

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

	if !p.peek(token.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStmt()

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

/*
func (p *Parser) parseCallArgs() []ast.Expr {
	args := []ast.Expr{}

	if p.isPeekToken(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpr(LOWEST))

	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpr(LOWEST))
	}

	if !p.peek(token.RPAREN) {
		return nil
	}

	return args
}
*/

func (p *Parser) GetErrors() []errs.ParserError {
	return p.errs
}

func (p *Parser) regPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) regInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekErr(t token.TokenType) {
    fmt.Println(p.curTok.Column , p.curTok.LineNo , p.curTok.Literal)
    fmt.Println(p.lx.GetLine(p.curTok.LineNo))
	expectedToken := t
	if len(t) > 1 {
		expectedToken = token.TokenType(token.HumanFriendly[string(t)])
	}
	//msg := fmt.Sprintf(errs.Errs["EXPECTED_GOT"], expectedToken, p.peekTok.Literal)
	newerr := errs.PeekError{Expected: expectedToken, Got: p.peekTok}
	p.errs = append(p.errs, &newerr)
	//p.errs = append(p.errs, errs.NewParserError(errs.EXPECTED_GOT, p.peekTok))
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.lx.NextToken()
}

func (p *Parser) ParseProg() *ast.Program {
	prog := &ast.Program{}
	prog.Stmts = []ast.Stmt{}

	for p.curTok.Type != token.EOF {

		//fmt.Println(p.curTok)
		stmt := p.parseStmt()

		if stmt != nil {
			prog.Stmts = append(prog.Stmts, stmt)
		}

		p.nextToken()
	}

	return prog
}

func (p *Parser) parseStmt() ast.Stmt {
	//fmt.Println(p.curTok.Type , p.peekTok)
	switch p.curTok.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.INCLUDE:
		return p.parseIncludeStmt()
	case token.COMMENT:
		return p.parseComment()
	default:
		return p.parseExprStmt()

	}
}

func (p *Parser) parseComment() ast.Stmt {

	return &ast.Comment{Token: p.curTok, Value: p.curTok.Literal}
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

func (p *Parser) parseIncludeStmt() *ast.IncludeStmt {
	stmt := &ast.IncludeStmt{Token: p.curTok}
	p.nextToken()

	stmt.Filename = p.parseExpr(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	log.Info(fmt.Sprintf("INCLUDE => FNAME=>%s || FNAME_TYPE=>%s", stmt.Filename, stmt))

	return stmt
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	//LET <IDENTIFIER> <EQUAL_SIGN> <EXPRESSION>
	stmt := &ast.LetStmt{Token: p.curTok}

	if !p.peek(token.IDENT) {
		return nil
	}

	stmt.Name = ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
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

func (p *Parser) noPrefixFunctionErr(t token.TokenType) {
	var msg errs.ParserError
    fmt.Println(string(p.lx.GetLine(p.curTok.LineNo)[p.curTok.Column]))
    //fmt.Println(p.curTok.Literal , " <-> " ,  p.lx.GetLine(p.curTok.LineNo))
	if t == token.FUNC {
		msg = &errs.NoEktiError{Type: t}
	} else {
		msg = &errs.NoPrefixSuffixError{Type: t}
	}
	p.errs = append(p.errs, msg)
}

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
		p.noPrefixFunctionErr(p.curTok.Type)
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

func (p *Parser) parseIdent() ast.Expr {
	log.Info("IDENT EXPR =>", p.curTok)
	return &ast.Identifier{
		Token: p.curTok,
		Value: p.curTok.Literal,
	}

}

func (p *Parser) parseBool() ast.Expr {
	log.Info("BOOL EXPR => ", p.curTok)
	return &ast.Boolean{Token: p.curTok, Value: p.isCurToken(token.TRUE)}
}

func (p *Parser) parseNumLit() ast.Expr {
	lit := &ast.NumberLit{Token: p.curTok}

	if number.IsFloat(p.curTok.Literal) {
		v, _ := new(big.Float).SetString(p.curTok.Literal)
		lit.Value = number.Number{Value: &number.FloatNumber{Value: *v}, IsInt: false}
		lit.IsInt = false
	} else {
		v, _ := new(big.Int).SetString(p.curTok.Literal, 10)
		lit.Value = number.Number{Value: &number.IntNumber{Value: *v}, IsInt: true}
		lit.IsInt = true
	}

	//if err != nil{
	//    return nil
	//}

	//lit.IsInt = value.IsInteger()
	//lit.Value = value

	return lit
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

func (p *Parser) parseIfExpr() ast.Expr {
	exp := &ast.IfExpr{Token: p.curTok}
	has_else := false
	if !p.peek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Cond = p.parseExpr(LOWEST)
	//fmt.Println(exp.Cond)

	if !p.peek(token.RPAREN) {
		return nil
	}
	// jodi (sotto) tahole { "hello" }
	if !p.peek(token.TAHOLE) {
		return nil
	}

	if !p.peek(token.LBRACE) {
		return nil
	}

	exp.TrueBlock = p.parseBlockStmt()

	if p.isPeekToken(token.ELSE) {
		p.nextToken()

		if !p.peek(token.LBRACE) {
			return nil
		}
		has_else = true
		exp.ElseBlock = p.parseBlockStmt()
	}

	if has_else {
		log.Info("IF ELSE Expr => ", exp.Cond, exp.TrueBlock.String(), exp.ElseBlock.String())
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

	if !p.peek(token.LBRACE) {
		return nil
	}

	exp.StmtBlock = p.parseBlockStmt()

	return exp

}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	bs := &ast.BlockStmt{Token: p.curTok}

	bs.Stmts = []ast.Stmt{}

	p.nextToken()

	for !p.isCurToken(token.RBRACE) && !p.isCurToken(token.EOF) {
		s := p.parseStmt()
		if s != nil {
			bs.Stmts = append(bs.Stmts, s)
		}
		p.nextToken()
	}
	//fmt.Println("BS=> " , bs)

	return bs
}

// Helper functions
func (p *Parser) isCurToken(t token.TokenType) bool {
	// check if current token type is `t`
	return p.curTok.Type == t
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
	// check if next token type is `t`
	return p.peekTok.Type == t
}

func (p *Parser) peek(t token.TokenType) bool {
	// checks if next token type is `t`
	// and if yes, then advance to next token
	if p.isPeekToken(t) {
		p.nextToken()
		return true
	} else {
		p.peekErr(t)
		return false
	}
}

// Check precedence of Peek Token
// (Token after current token)
func (p *Parser) peekPrec() int {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}

	return LOWEST
}

// Check precedence of Current Token
func (p *Parser) curPrec() int {
	if p, ok := precedences[p.curTok.Type]; ok {
		return p
	}
	return LOWEST
}
