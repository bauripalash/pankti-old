package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	LineNo  int
	Column  int
}

const (

	//Symbols
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	PLUS    = "+"

	STRING = "STRING"
	// Identifier token
	IDENT = "IDENT"

	//Left Square Bracket `[`
	LS_BRACKET = "["
	//Rigt Square Bracket `]`
	RS_BRACKET = "]"

	COLON = ":"

	// integer
	INT = "INT"

	FLOAT = "FLOAT"

    NUM = "NUMBER"

	//Equal = sign; for assignment
	EQ = "="

	EQEQ   = "=="
	NOT_EQ = "!="
	MUL    = "*"
	DIV    = "/"
	MINUS  = "-"

	//Bang or `!`
	EXC = "!"

	LT        = "<"
	LTE       = "<="
	GT        = ">"
	GTE       = ">="
	SEMICOLON = ";"
	COMMA     = ","
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	HOLO      = "HOLO"
	EKTI      = "EKTI"
	TAHOLE    = "TAHOLE"

	//Keywords

	FUNC   = "FUNCTION"
	LET    = "LET"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
    WHILE  = "WHILE"
)

var HumanFriendly = map[string]string{

	IDENT:  "নাম",
	FUNC:   "kaj",
	LET:    "dhori",
	TRUE:   "sotti",
	FALSE:  "mittha",
	IF:     "jodi",
	ELSE:   "nahole",
	RETURN: "ferau",
	HOLO:   "holo",
	EKTI:   "ekti",
	TAHOLE: "tahole",
    WHILE: "jotokhon",
}

var Keywords = map[string]TokenType{

	"কাজ":    FUNC,
	"kaj":    FUNC,
	"fn":     FUNC,
	"ধরি":    LET,
	"dhori":  LET,
	"let":    LET,
	"সত্য":   TRUE,
	"sotto":  TRUE,
	"মিথ্যা": FALSE,
	"mittha": FALSE,
	"যদি":    IF,
	"jodi":   IF,
	"নাহলে":  ELSE,
	"nahole": ELSE,
	"ফেরাও":  RETURN,
	"ferau":  RETURN,
	"হল":     HOLO,
	"holo":   HOLO,
	"একটি":   EKTI,
	"ekti":   EKTI,
	"তাহলে":  TAHOLE,
	"tahole": TAHOLE,
    "jotokhon" : WHILE,
    "while" : WHILE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}

	return IDENT
}
