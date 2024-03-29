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

	NUM     = "NUMBER"
	INCLUDE = "INCLUDE"

	COMMENT = "COMMENT"

	//Equal = sign; for assignment
	EQ = "="

	EQEQ   = "=="
	NOT_EQ = "!="
	MUL    = "*"
	DIV    = "/"
	MINUS  = "-"
	MOD    = "%"

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
	AND       = "AND"
	OR        = "OR"

	//Keywords

	FUNC   = "FUNCTION"
	LET    = "LET"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	WHILE  = "WHILE"
	SHOW   = "SHOW"
	END    = "END"
	DUMMY  = "DUMMY"
	BREAK  = "BREAK"
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
	WHILE:  "jotokhon",
	SHOW:   "dekhau",
	END:    "end",
	BREAK:  "break",
	AND:    "and",
	OR:     "or",
}

var Keywords = map[string]TokenType{

	"কাজ":      FUNC,
	"kaj":      FUNC,
	"fn":       FUNC,
	"ধরি":      LET,
	"dhori":    LET,
	"let":      LET,
	"সত্য":     TRUE,
	"sotto":    TRUE,
	"মিথ্যা":   FALSE,
	"mittha":   FALSE,
	"যদি":      IF,
	"jodi":     IF,
	"নাহলে":    ELSE,
	"nahole":   ELSE,
	"else":     ELSE,
	"ফেরাও":    RETURN,
	"ferau":    RETURN,
	"ferao":    RETURN,
	"হল":       HOLO,
	"holo":     HOLO,
	"একটি":     EKTI,
	"ekti":     EKTI,
	"তাহলে":    TAHOLE,
	"tahole":   TAHOLE,
	"jotokhon": WHILE,
	"while":    WHILE,
	"যতক্ষণ":   WHILE,
	"include":  INCLUDE,
	"anoyon":   INCLUDE,
	"আনয়ন":     INCLUDE,
	"dekhau":   SHOW,
	"dekhao":   SHOW,
	"show":     SHOW,
	"দেখাও":    SHOW,
	"sesh":     END,
	"end":      END,
	"শেষ":      END,
	"break":    BREAK,
	"সমাপ্ত":   BREAK,
	"somapto":  BREAK,
	"ebong":    AND,
	"and":      AND,
	"এবং":      AND,
	"ba":       OR,
	"or":       OR,
	"বা":       OR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}

	return IDENT
}
