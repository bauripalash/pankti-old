package errs

import (
	"fmt"

	"go.cs.palashbauri.in/pankti/token"
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
	ErrLine  string
}

func (pe *PeekError) GetMsg() string { return Errs[EXPECTED_GOT] }

func (pe *PeekError) GetToken() token.Token { return pe.Got }

func (pe *PeekError) String() string {
	return pe.ErrLine + "\n" + fmt.Sprintf(
		pe.GetMsg(),
		pe.Expected,
		pe.GetToken().Literal,
	)
}

type NoPrefixSuffixError struct {
	Token   token.Token
	ErrLine string
	//Type token.TokenType
}

func (spe *NoPrefixSuffixError) GetMsg() string {
	return Errs[NO_PREFIX_SUFFIX_FN]
}

func (spe *NoPrefixSuffixError) GetToken() token.Token {
	return spe.Token
}

func (spe *NoPrefixSuffixError) String() string {
	//fmt.Println()
	return spe.ErrLine + "\n" + fmt.Sprintf(
		spe.GetMsg(),
		spe.Token.Literal,
	)

}

type NoEktiError struct {
	Type    token.TokenType
	ErrLine string
}

func (nee *NoEktiError) GetMsg() string { return Errs[NO_EKTI_BEFORE_FN] }

func (nee *NoEktiError) GetToken() token.Token { return token.Token{} }

func (nee *NoEktiError) String() string {
	return nee.ErrLine + "\n" + fmt.Sprintf(
		Errs[NO_EKTI_BEFORE_FN],
		nee.Type,
	)
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

	"NO_EKTI_BEFORE_FN":              "`কাজ`-এর আগে 'ekti' বা 'একটি' পাওয়া উচিত ছিল %s",
	"EXPECTED_GOT":                   "এখানে `%s` পাওয়া উচিত ছিল কিন্তু `%s` পাওয়া গেল",
	"NO_PREFIX_SUFFIX_FN":            "এটা %s নিয়ে কী করা উচিত আমি জানিনা",
	"INT_PARSE_ERR":                  "%s - এই এটা তো একটা সংখ্যা নয়",
	"FUN_CALL_NOT_ENOUGH_ARGS":       "এই '%s' কাজের জন্য %dটি চল রাশির প্রয়োজন কিন্তু পাওয়া গেলো %dটি",
	"NOT_ALL_ARE_INT":                "এই '%s' কাজের সমস্ত জন্য দেওয়া সব চলরাশি গুলিকে সংখ্যা হতে হবে",
	"NOT_ALL_ARE_LIST":               "এই কাজের জন্য প্রদত্ত সমস্ত চলরাশি গুলিকে 'তালিকা' হতে হবে।",
	"INDEX_MUST_BE_NUMBER":           "এই কাজের জন্য সূচকটিকে সংখ্যা হতে হবে।",
	"INDEX_OUT_RANGE":                "এই সূচকটি তালিকার আয়তনের থেকেও বড়।",
	"FILENAME_MUST_BE_STRING":        "এখানে ফাইলের নাম একটি স্ট্রিং বা 'লেখা নাম' হতে হবে।",
	"FAILED_TO_READ_FILE":            "এই ফাইলটি পড়া গেলো না।",
	"FAILED_TO_CREATE":               "এই ফাইলটি '%s' তৈরি করা গেল না!",
	"FAILED_TO_CLOSE_FILE":           "এই ফাইলটি '%s' তৈরি করে খোলার পর আর বন্ধ করা গেল না।",
	"FAILED_TO_WRITE_FILE":           "এই ফাইলটিতে '%s' লেখা গেল না।",
	"FILE_PATH_MUST_BE_STRING":       "ফাইলের ঠিকানা বা প্যাথ একটি স্ট্রিং/'লেখা' হতে হবে",
	"FILE_NOT_EXIST":                 "এই ফাইলটির অস্তিত্ব খুঁজে পাওয়া গেল না।",
	"NEW_FILENAME_MUST_BE_STRING":    "নতুন ফাইলের নাম স্ট্রিং/'লেখা' হতে হবে",
	"RENAME_FAILED":                  "ফাইলটির নাম পরিবর্তন করা গেল না।",
	"DELETE_FAILED":                  "ফাইলটি মুছে ফেলা বা ডিলিট করা গেল না।",
	"DATA_MUST_BE_STRING":            "ফাইলের লেখার জন্য দেওয়া তথ্যকে একটি স্ট্রিং বা 'লেখা' হতে হবে।",
	"FAILED_OPEN_FILE":               "ফাইলটি খোলা গেল না।",
	"FAILED_TO_WRITE_DATA":           "ফাইলে প্রদত্ত তথ্য সংরক্ষিত করা গেল না",
	"TARGET_NO_DIR":                  "প্রদত্ত ঠিকানাটি কোন ফোল্ডার/ডাইরেক্টরি কে নির্দেশ করে না",
	"TARGET_IS_DIR":                  "প্রদত্ত ঠিকানাটি কোন ফোল্ডার/ডাইরেক্টরি কে নির্দেশ করে",
	"DIR_LIST_FAILED":                "ফোল্ডার/ডাইরেক্টরির ফাইলগুলির সূচি তৈরি করা গেল না।",
	"ARG_DECIMAL_PARSE_FAILED":       "কাজের জন্য প্রদত্ত চল রাশিগুলি দশমিক সংখ্যা হিসাবে গ্রহণ করা গেল না।",
	"GCD_ALL_INT":                    "গসাগুর জন্য প্রদত্ত সমস্ত সংখ্যাগুলিকে পূর্ণসংখ্যা/Integer হতে হবে",
	"SUM_ONLY_LISTS":                 "যোগফল শুধুমাত্র সংখ্যাযুক্ত তালিকারই বার করা সম্ভব",
	"SUM_ARRAY_ALL_NUM":              "যোগফল বার করার জন্য তালিকার সমস্ত সদস্যকে সংখ্যা হতে হবে।",
	"TEMPLATE_NOT_ALL_INT":           "এই '%s' কাজটি করার জন্য সমস্ত প্রদত্ত চলরাশিগুলিকে %s হতে হবে।",
	"NOT_ALL_STRING":                 "এই কাজটি করার জন্য সমস্ত প্রদত্ত চলরাশিগুলিকে স্ট্রিং/'লেখা' হতে হবে।",
	"TEMPLATE_NOT_ONE_TEMPALTE":      "এই কাজটি '%s' করার জন্য প্রদত্ত চলরাশিকে '%s' হতে হবে।",
	"STRING_NUM_PARSE_FAIL_TEMPLATE": "প্রদত্ত চলরাশিকে স্ট্রিং/'লেখা' হিসাবে গ্রহণ করা গেল না।",
	"CANNOT_PARSE_AS_NUM":            "প্রদত্ত চলরাশিকে সংখ্যাতে পরিণত করা যাবে না।",
	"CANNOT_PARSE_STRING_AS_NUM":     "প্রদত্ত স্ট্রিং/'লেখা'কে সংখ্যাতে পরিণত করা যাবে না।",
	"STDIN_READ_FAILED":              "Stdin থেকে তথ্য পড়া গেল না।",
}
