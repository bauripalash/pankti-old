package stdlib

import (
	"fmt"
	"time"
	"vabna/object"
)

func UnixTimeFunc(args []object.Obj) object.Obj {

	return &object.String{Value: fmt.Sprintf("%d", time.Now().Unix())}
}
