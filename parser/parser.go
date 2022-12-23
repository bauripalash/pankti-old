package parser

import (
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/token"
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
	p.regPrefix(token.INCLUDE, p.parseIncludeExpr)

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

func (p *Parser) regPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) regInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

		//		if stmt != nil {
		prog.Stmts = append(prog.Stmts, stmt)
		//		}

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
	//case token.INCLUDE:
	//	return p.parseIncludeStmt()
	case token.COMMENT:
		return p.parseComment()
	case token.SHOW:
		return p.parseShowStmt()
	default:
		return p.parseExprStmt()

	}
}

func (p *Parser) parseComment() ast.Stmt {

	return &ast.Comment{
		Token: p.curTok,
		Value: p.curTok.Literal,
	}
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
