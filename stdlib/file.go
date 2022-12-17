package stdlib

import (
	"os"

	"go.cs.palashbauri.in/pankti/object"
)

func getStringFromArgs(args []object.Obj) (string, bool) { // Value, isOkay
	fileArg := args[0]

	if fileArg.Type() == object.STRING_OBJ {
		return fileArg.(*object.String).Value, true
	}
	return "", false

}

func ReadFile(args []object.Obj) object.Obj {
	filename, isOkay := getStringFromArgs(args)

	if !isOkay {
		return &object.Error{Msg: "Filename must be string"}
	}
	d, err := os.ReadFile(filename)
	if err != nil {
		return &object.Error{Msg: "Failed to read file"}
	}

	return &object.String{Value: string(d)}
}

func FileDirExists(args []object.Obj) object.Obj {
	result := false
	filename, isOkay := getStringFromArgs(args)

	if !isOkay {
		return &object.Error{Msg: "Filename must be string"}
	}

	if _, err := os.Stat(filename); err == nil {
		result = true
	}
	return &object.Boolean{Value: result}
}
