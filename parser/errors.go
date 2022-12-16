package parser

import (
	"strconv"

	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/token"
)

func (p *Parser) GetErrors() []errs.ParserError {
	return p.errs
}

func (p *Parser) peekErr(t token.TokenType) {
	expectedToken := t
	if len(t) > 1 {
		expectedToken = token.TokenType(token.HumanFriendly[string(t)])
	}
	newerr := errs.PeekError{
		Expected: expectedToken,
		Got:      p.peekTok,
		ErrLine:  MakeErrorLine(p.curTok, p.lx.GetLine(p.curTok.LineNo)),
	}
	p.errs = append(p.errs, &newerr)
}

func MakeErrorLine(t token.Token, line string) string {
	//    fmt.Println(t.LineNo , line)
	Lindex := t.Column - 1

	RIndex := t.Column + len(t.Literal) - 1

	if len(t.Literal) <= 1 {
		RIndex = Lindex + 1
	}

	newL := line[:RIndex] + " <-- " + line[RIndex:]
	newLine := newL[:Lindex] + " --> " + newL[Lindex:]

	return strconv.Itoa(t.LineNo) + "| " + newLine
}

func (p *Parser) noPrefixFunctionErr(t token.Token) {
	var msg errs.ParserError

	if t.Type == token.FUNC {
		msg = &errs.NoEktiError{
			Type:    t.Type,
			ErrLine: MakeErrorLine(t, p.lx.GetLine(t.LineNo)),
		}
	} else {
		msg = &errs.NoPrefixSuffixError{
			Token:   p.curTok,
			ErrLine: MakeErrorLine(t, p.lx.GetLine(t.LineNo)),
		}

	}
	p.errs = append(p.errs, msg)
}
