package stdlib

import (
	"runtime"

	ob "go.cs.palashbauri.in/pankti/object"
)

// Operating System Related Queries
func GetOS(_ []ob.Obj) ob.Obj {
	return &ob.String{Value: runtime.GOOS}
}

func GetArch(_ []ob.Obj) ob.Obj {
	return &ob.String{Value: runtime.GOARCH}
}
