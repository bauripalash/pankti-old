package lexer

import (
	"strings"

	"go.cs.palashbauri.in/pankti/token"
)

type Lexer struct {
	input   []rune
	pos     int
	readPos int
	ch      rune
	line    int
	column  int
}

func (l *Lexer) AtEOF() bool {

	return l.pos >= len(l.input)

}

/*
func getLen(inp string) int {

	return utf8.RuneCountInString(inp)

}
*/

func NewLexer(input string) Lexer {
	lexer := Lexer{input: []rune(input), line: 1, column: 0}
	lexer.readChar()
	return lexer
}

func (l *Lexer) readChar() {
	//Advances lexer

	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	//fmt.Printf("<-> %c >> %d >>  %d >> %d\n", l.ch, len(string(l.ch)), l.pos, l.readPos)
	l.pos = l.readPos

	l.readPos++
	l.column++
}

func (l *Lexer) GetLine(ln int) string {
	a := strings.Split(string(l.input), "\n")

	for idx, l := range a {
		if idx == (ln - 1) {
			return l
		}
		//fmt.Println(idx , l)
	}
	//b:= strings.Split(string(l.input) , "\n")

	//if len(a) > len(b) {
	//    fmt.Println("a->" , a[ln] , "<=>" )
	//}else{
	// fmt.Println("b->", a , "<=>" , len(a))
	//}

	//fmt.Println(a)
	return ""
}

func (l *Lexer) NextToken() token.Token {
	// Get next token

	var tk token.Token
	l.eatWhitespace()
	switch l.ch {

	case '+':
		tk = NewToken(token.PLUS, l.ch, l.line, l.column)
	case '-':
		tk = NewToken(token.MINUS, l.ch, l.line, l.column)
	case '*':
		tk = NewToken(token.MUL, l.ch, l.line, l.column)
	case '/':
		tk = NewToken(token.DIV, l.ch, l.line, l.column)
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tk = token.Token{
				Type:    token.EQEQ,
				Literal: lit,
				LineNo:  l.line,
				Column:  l.column,
			}
		} else {
			tk = NewToken(token.EQ, l.ch, l.line, l.column)
		}
	case ';':
		tk = NewToken(token.SEMICOLON, l.ch, l.line, l.column)
	case ',':
		tk = NewToken(token.COMMA, l.ch, l.line, l.column)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tk = token.Token{
				Type:    token.LTE,
				Literal: lit,
				LineNo:  l.line,
				Column:  l.column,
			}
		} else {
			tk = NewToken(token.LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tk = token.Token{
				Type:    token.GTE,
				Literal: lit,
				LineNo:  l.line,
				Column:  l.column,
			}
			//fmt.Println(tk)
		} else {
			tk = NewToken(token.GT, l.ch, l.line, l.column)
		}
	case '(':
		tk = NewToken(token.LPAREN, l.ch, l.line, l.column)
	case ')':
		tk = NewToken(token.RPAREN, l.ch, l.line, l.column)
	case '{':
		tk = NewToken(token.LBRACE, l.ch, l.line, l.column)
	case '}':
		tk = NewToken(token.RBRACE, l.ch, l.line, l.column)
	case '%':
		tk = NewToken(token.MOD, l.ch, l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tk = token.Token{
				Type:    token.NOT_EQ,
				Literal: lit,
				LineNo:  l.line,
				Column:  l.column,
			}
		} else {
			tk = NewToken(token.EXC, l.ch, l.line, l.column)
		}
	case '"':
		tk.Type = token.STRING
		tk.LineNo = l.line
		tk.Column = l.column
		tk.Literal = l.readString()
	case '[':
		tk = NewToken(token.LS_BRACKET, l.ch, l.line, l.column)
	case ']':
		tk = NewToken(token.RS_BRACKET, l.ch, l.line, l.column)
	case ':':
		tk = NewToken(token.COLON, l.ch, l.line, l.column)
		//TODO: case '#':

		//TODO: Comments l.eatSingleLineComment()
	case '#':
		pos := l.pos + 1
		col := l.column
		lin := l.line

		for {

			l.readChar()

			if l.ch == '\n' || l.ch == '\r' || l.ch == 0 {
				break
			}

		}

		tk = token.Token{
			Type:    token.COMMENT,
			Literal: string(l.input[pos:l.pos]),
			LineNo:  lin,
			Column:  col,
		}

		//l.eatWhitespace()
		//log.Print("->")
		//log.Print(l.ch)
		//log.Println("<-")

		//l.eatWhitespace()

	case 0:
		tk.Column = l.column
		tk.LineNo = l.line
		tk.Literal = ""
		tk.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tk.LineNo = l.line
			tk.Column = l.column
			lit, _ := l.readIdent()
			tk.Literal = lit
			tk.Type = token.LookupIdent(tk.Literal)

			return tk
		} else if isDigit(l.ch) {
			tk.LineNo = l.line
			tk.Column = l.column
			lit, _ := l.readNum()

			//fmt.Println(lit)
			tk.Literal = lit
			tk.Type = token.NUM

			return tk
		} else {
			tk = NewToken(token.ILLEGAL, l.ch, l.line, l.column)
		}

	}

	l.readChar()
	return tk

}

/*
func (l *Lexer) eatSingleLineComment(){

    for l.ch != '\n' && !l.AtEOF(){
        l.readChar()
        //l.eatWhitespace()
    }


    //l.eatWhitespace()

    log.Println("->" + string(l.ch) + "<-")
}
*/

func (l *Lexer) readString() string {
	pos := l.pos + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	//fmt.Println(l.input[pos:l.pos])
	return string(l.input[pos:l.pos])
}

func (l *Lexer) eatWhitespace() {
	for l.ch == rune(' ') || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
		//log.Println(l.ch)

	}
}

func NewToken(
	tokType token.TokenType,
	ch rune,
	line int,
	col int,
) token.Token {

	return token.Token{
		Type:    tokType,
		Literal: string(ch),
		LineNo:  line,
		Column:  col,
	}

}

func (l *Lexer) readIdent() (string, bool) {

	pos := l.pos
	isModId := false
	for isLetter(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		isModId = true
		l.readChar()
		for isLetter(l.ch) {
			l.readChar()
		}
	}
	return string(l.input[pos:l.pos]), isModId

}

func (l *Lexer) readNum() (string, bool) {
	pos := l.pos
	isFloat := false

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		isFloat = true
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	//parseBengaliNum(l.input[pos:l.pos])

	return string(parseBengaliNum(l.input[pos:l.pos])), isFloat
}

func parseBengaliNum(inp []rune) []rune {

	var result []rune

	for _, item := range inp {

		//fmt.Println(item)

		switch item {
		case '০':
			result = append(result, '0')
		case '১':
			result = append(result, '1')
		case '২':
			result = append(result, '2')
		case '৩':
			result = append(result, '3')
		case '৪':
			result = append(result, '4')
		case '৫':
			result = append(result, '5')
		case '৬':
			result = append(result, '6')
		case '৭':
			result = append(result, '7')
		case '৮':
			result = append(result, '8')
		case '৯':
			result = append(result, '9')
		default:
			result = append(result, item)
		}
	}

	//fmt.Println(string(result))
	return result
}

func (l *Lexer) peekChar() rune {

	if l.readPos >= len(l.input) {
		return 0
	} else {
		return []rune(l.input)[l.readPos]
	}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' ||
		'ঀ' <= ch && ch <= 'ৡ' ||
		'ৰ' <= ch && ch <= '৽'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || '০' <= ch && ch <= '৯'
}
