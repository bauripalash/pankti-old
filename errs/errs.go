package errs

import (
	"fmt"
	"vabna/token"
)

const (
	NO_EKTI_BEFORE_FN   = "NO_EKTI_BEFORE_FN"
	EXPECTED_GOT        = "EXPECTED_GOT"
	NO_PREFIX_SUFFIX_FN = "NO_PREFIX_SUFFIX_FN"
	INT_PARSE_ERR       = "INT_PARSE_ERR"
)

type ParserError interface {
	GetMsg() string
	GetToken() token.Token
	String() string
}

type PeekError struct {
	Msg      string
	Expected token.TokenType
	Got      token.Token
}

func (pe *PeekError) GetMsg() string { return Errs[EXPECTED_GOT] }

func (pe *PeekError) GetToken() token.Token { return pe.Got }

func (pe *PeekError) String() string {
	return fmt.Sprintf(pe.GetMsg(), pe.Expected, pe.GetToken().Literal)
}

type NoPrefixSuffixError struct {
	Type token.TokenType
}

func (spe *NoPrefixSuffixError) GetMsg() string {
	return Errs[NO_PREFIX_SUFFIX_FN]
}

func (spe *NoPrefixSuffixError) GetToken() token.Token {
	return token.Token{}
}

func (spe *NoPrefixSuffixError) String() string {

	return fmt.Sprintf(spe.GetMsg(), spe.Type)

}

type NoEktiError struct {
	Type token.TokenType
}

func (nee *NoEktiError) GetMsg() string { return Errs[NO_EKTI_BEFORE_FN] }

func (nee *NoEktiError) GetToken() token.Token { return token.Token{} }

func (nee *NoEktiError) String() string {
	return fmt.Sprintf(Errs[NO_EKTI_BEFORE_FN], nee.Type)
}

type IntegerParseError struct {
	Token token.Token
}

func (ipe *IntegerParseError) GetMsg() string { return Errs[INT_PARSE_ERR] }

func (ipe *IntegerParseError) GetToken() token.Token { return ipe.Token }

func (ipe *IntegerParseError) String() string {
	return fmt.Sprintf(ipe.GetMsg(), ipe.GetToken())
}

var Errs = map[string]string{

	"NO_EKTI_BEFORE_FN":   "`কাজ`-এর আগে 'ekti' বা 'একটি' পাওয়া উচিত ছিল %s",
	"EXPECTED_GOT":        "এখানে `%s` পাওয়া উচিত ছিল কিন্তু `%s` পাওয়া গেল",
	"NO_PREFIX_SUFFIX_FN": "এটা %s নিয়ে কী করা উচিত আমি জানিনা",
	"INT_PARSE_ERR":       "%s - এই এটা তো একটা সংখ্যা নয়",
}
