package ast

import (
	"bytes"
	"go.cs.palashbauri.in/pankti/token"
	"strings"
)

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
