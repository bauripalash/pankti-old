package lexer

import (
	"testing"

	"go.cs.palashbauri.in/pankti/token"
)

func TestNextToken(t *testing.T) {

	inp := `let name = 100;
    let hello = "hello";
    dhori nam = "Palash Bauri";
    dhori nach = ekti kaj(lok) {
        dekhau(lok + " is dancing");
    }
    dhori sonkhya = ১০০;
    ১0০৯
    `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "name"},
		{token.EQ, "="},
		{token.NUM, "100"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "hello"},
		{token.EQ, "="},
		{token.STRING, "hello"},
		{token.SEMICOLON, ";"},
		{token.LET, "dhori"},
		{token.IDENT, "nam"},
		{token.EQ, "="},
		{token.STRING, "Palash Bauri"},
		{token.SEMICOLON, ";"},
		{token.LET, "dhori"},
		{token.IDENT, "nach"},
		{token.EQ, "="},
		{token.EKTI, "ekti"},
		{token.FUNC, "kaj"},
		{token.LPAREN, "("},
		{token.IDENT, "lok"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "dekhau"},
		{token.LPAREN, "("},
		{token.IDENT, "lok"},
		{token.PLUS, "+"},
		{token.STRING, " is dancing"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "dhori"},
		{token.IDENT, "sonkhya"},
		{token.EQ, "="},
		{token.NUM, "100"},
		{token.SEMICOLON, ";"},
		{token.NUM, "1009"},
	}

	l := NewLexer(inp)

	for i, tt := range tests {
		tk := l.NextToken()

		if tk.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] -> TokenType wrong -> Expected=%q, Got=%q",
				i,
				tt.expectedType,
				tk.Type,
			)
		}

		if tk.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] -> Literal wrong -> Expected=%q, Got=%q",
				i,
				tt.expectedLiteral,
				tk.Literal,
			)
		}
	}

}
