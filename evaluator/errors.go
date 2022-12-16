package evaluator

import (
	"fmt"
	"strconv"
	"strings"

	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/token"
)

type ErrorHelper struct {
	Source string
}

func (e *ErrorHelper) GetLine(t token.Token) string {

	return strings.Split(e.Source, "\n")[t.LineNo-1]
}

func (e *ErrorHelper) MakeErrorLine(t token.Token, showHint bool) string {

	if t.LineNo <= 0 { //if the token is virtual token; line can be zero
		return ""
	}

	//newLine := e.Source
	var newLine string
	xLine := e.GetLine(t)
	if showHint {

		//fmt.Println(xLine)

		Lindex := t.Column - 1
		if Lindex < 0 { //In case of a virtual token
			Lindex = 0
		}

		RIndex := t.Column + len(t.Literal) - 1

		if len(t.Literal) <= 1 {
			RIndex = Lindex + 1
		}

		newL := xLine[:RIndex] + " <-- " + xLine[RIndex:]
		newLine = newL[:Lindex] + " --> " + newL[Lindex:]

		return strconv.Itoa(t.LineNo) + "| " + newLine
	}
	return strconv.Itoa(t.LineNo) + "| " + xLine
}

func NewErr(
	token token.Token,
	eh *ErrorHelper,
	showHint bool,
	format string,
	a ...interface{},
) *object.Error {

	errMsg := eh.MakeErrorLine(
		token,
		showHint,
	) + "\n" + fmt.Sprintf(
		format,
		a...)
	return &object.Error{Msg: errMsg}
}

func isErr(obj object.Obj) bool {
	if obj != nil {
		return obj.Type() == object.ERR_OBJ
	}

	return false
}

func NewBareErr(format string, a ...interface{}) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}
