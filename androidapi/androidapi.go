package androidapi

import (
	"strings"

	"bauri.palash/pankti/evaluator"
	"bauri.palash/pankti/lexer"
	"bauri.palash/pankti/object"
	"bauri.palash/pankti/parser"
)




func DoParse(i string) string{
    l := lexer.NewLexer(i)
    p := parser.NewParser(&l)
    prog := p.ParseProg()

    if len(p.GetErrors()) >= 1{
        tempErrs := []string{}

        for _,item := range p.GetErrors(){
            tempErrs = append(tempErrs, item.String())
        }

        return strings.Join(tempErrs , " \n")
    }

    env := object.NewEnv()
    eh := evaluator.ErrorHelper{ Source: i }
    evd := evaluator.Eval(prog , env , eh)

    if evd != nil{
        return evd.Inspect()
    }else{
        return ""
    }
}
