package stdlib

import (
	"fmt"
	"go.cs.palashbauri.in/pankti/object"
	"time"
)

func UnixTimeFunc(args []object.Obj) object.Obj {

	return &object.String{Value: fmt.Sprintf("%d", time.Now().Unix())}
}
