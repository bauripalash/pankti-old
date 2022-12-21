package stdlib

import (
	"os/user"
	"runtime"

	ob "go.cs.palashbauri.in/pankti/object"
)

// Operating System Related Queries
func GetOS() ob.Obj {
	return &ob.String{Value: runtime.GOOS}
}

func GetArch() ob.Obj {
	return &ob.String{Value: runtime.GOARCH}
}

func GetUserName() ob.Obj {
	u, err := user.Current()
	if err != nil {
		return &ob.String{Value: ""}
	}

	return &ob.String{Value: u.Username}

}

func GetUserHomeDir() ob.Obj {
	u, err := user.Current()
	if err != nil {
		return &ob.String{Value: ""}
	}

	return &ob.String{Value: u.HomeDir}
}
