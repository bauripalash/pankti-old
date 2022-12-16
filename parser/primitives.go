package parser

import (
	"math/big"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) parseStringLit() ast.Expr {

	//fmt.Println(p.curTok)
	return &ast.StringLit{Token: p.curTok, Value: p.curTok.Literal}
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
	return &ast.Boolean{
		Token: p.curTok,
		Value: p.isCurToken(token.TRUE),
	}
}

func (p *Parser) parseNumLit() ast.Expr {

	lit := &ast.NumberLit{Token: p.curTok}

	if number.IsFloat(p.curTok.Literal) {
		v, _ := new(big.Float).SetString(p.curTok.Literal)
		lit.Value = number.Number{
			Value: &number.FloatNumber{Value: *v},
			IsInt: false,
		}
		lit.IsInt = false
	} else {
		v, _ := new(big.Int).SetString(p.curTok.Literal, 10)
		lit.Value = number.Number{
			Value: &number.IntNumber{Value: *v},
			IsInt: true,
		}
		lit.IsInt = true
	}

	//if err != nil{
	//    return nil
	//}

	//lit.IsInt = value.IsInteger()
	//lit.Value = value

	return lit
}
