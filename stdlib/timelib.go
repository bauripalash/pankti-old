package stdlib

import (
	"fmt"
	"go.cs.palashbauri.in/pankti/object"
	"time"
)

func UnixTimeFunc(args []object.Obj) object.Obj {

	return &object.String{Value: fmt.Sprintf("%d", time.Now().Unix())}
}

func UtcDateISO(args []object.Obj) object.Obj {
	return &object.String{Value: fmt.Sprintf("%s", time.Now().Format(time.RFC3339))}
}
