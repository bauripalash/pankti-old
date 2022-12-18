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

func CreateEmptyFile(args []object.Obj) object.Obj {
	filename, isOkay := getStringFromArgs(args)

	if !isOkay {
		return &object.Error{Msg: "Filename must be string"}
	}

	if _, err := os.Stat(filename); err == nil {
		return &object.Boolean{Value: false}
	}

	f, err := os.Create(filename)
	if err != nil {
		return &object.Error{Msg: "Filed to Create File " + filename}
	}

	if err := f.Close(); err == nil {
		return &object.Error{Msg: "Error while Creating File " + filename + ";" + "Failed to close the file"}
	} else {
		return &object.Boolean{Value: true}
	}

}

func WriteToFile(args []object.Obj) object.Obj {
	filename, isOkay := getStringFromArgs(args)
	data := args[1]
	if !isOkay {
		return &object.Error{Msg: "Filename must be string"}
	}

	err := os.WriteFile(filename, []byte(data.Inspect()), 0644)

	if err != nil {
		return &object.Error{Msg: "Failed to write to file " + filename}
	}

	return &object.Boolean{Value: true}

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
