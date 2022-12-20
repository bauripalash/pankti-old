package object

import (
	"fmt"
	"strconv"
	"strings"

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
	//var newLine string
	xLine := e.GetLine(t)
	//fmt.Println(xLine)
	if showHint {

		x := []rune(xLine)
		var li int
		var ri int

		if temp := t.Column + len([]rune(t.Literal)); temp < len(x) {
			ri = temp
		} else {
			ri = len(x) - 1
		}

		//fmt.Println(xLine)

		//Lindex := t.Column - 1
		//fmt.Println(t.Column)
		if li < 0 { //In case of a virtual token
			li = 0
			//return ""
		}

		//RIndex := t.Column + len(t.Literal) - 1

		if len(t.Literal) <= 1 {
			ri = li + 1
		}

		li = t.Column

		temp_x := string(x[:ri]) + " <-- " + string(x[ri:])
		y := []rune(temp_x)
		temp_y := string(y[:li]) + " --> " + string(y[li:])
		//fmt.Println(temp_y)

		//newL := xLine[:RIndex] + " <-- " + xLine[RIndex:]
		//newLine = newL[:Lindex] + " --> " + newL[Lindex:]

		return strconv.Itoa(t.LineNo) + "| " + temp_y
	}
	return strconv.Itoa(t.LineNo) + "| " + xLine
}

func NewErr(
	token token.Token,
	eh *ErrorHelper,
	showHint bool,
	format string,
	a ...interface{},
) *Error {

	errMsg := eh.MakeErrorLine(
		token,
		showHint,
	) + "\n" + fmt.Sprintf(
		format,
		a...)
	return &Error{Msg: errMsg}
}

func IsErr(obj Obj) bool {
	if obj != nil {
		return obj.Type() == ERR_OBJ
	}

	return false
}

func NewBareErr(format string, a ...interface{}) Obj {
	return &Error{Msg: fmt.Sprintf(format, a...)}
}
