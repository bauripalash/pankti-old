package stdlib

import (
	"fmt"
	"time"

	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/object"
)

func UnixTimeFunc(args []object.Obj) object.Obj {

	return &object.String{Value: fmt.Sprintf("%d", time.Now().Unix())}
}

func UtcDateISO(args []object.Obj) object.Obj {
	return &object.String{
		Value: fmt.Sprintf("%s", time.Now().Format(time.RFC3339)),
	}
}

func TimeNow() object.Obj {
	return &object.String{
		Value: time.Now().Local().Format(time.Kitchen),
	}
}

func DateNow() object.Obj {
	return &object.String{
		Value: time.Now().Local().Format("02-01-2006"),
	}
}

// %y%m%d -> 2006 - 01 - 02
func dateFormat(d string) string {
	//TODO: Fix
	return "02-01-2006"
}

func FormatTimeLocal(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	f, ok := getStringFromArgs(args[0])

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	return &object.String{
		Value: time.Now().Local().Format(dateFormat(f)),
	}
}

func FormatTimeUTC(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	f, ok := getStringFromArgs(args[0])

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	return &object.String{
		Value: time.Now().UTC().Format(dateFormat(f)),
	}
}
