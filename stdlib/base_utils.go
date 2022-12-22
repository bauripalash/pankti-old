package stdlib

import (
	"bufio"
	"fmt"
	"os"

	"go.cs.palashbauri.in/pankti/constants"
	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/object"
)

func ReturnErrorString(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	msg, ok := getStringFromArgs(args[0])

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	return object.NewBareErr(msg)
}

func ReadLine(eh *object.ErrorHelper, args []object.Obj) object.Obj {

	msg, ok := getStringFromArgs(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(msg)
	text, err := reader.ReadString('\n')

	if err != nil {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["STDIN_READ_FAILED"])
	}

	return &object.String{Value: text}
}

func GetType(_ *object.ErrorHelper, args []object.Obj) object.Obj {
	t, ok := constants.TypeNames[string(args[0].Type())]
	if !ok {
		t = constants.UNKNOWN
	}
	return &object.String{Value: t}
}
