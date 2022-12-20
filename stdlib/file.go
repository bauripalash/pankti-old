package stdlib

import (
	"os"

	"go.cs.palashbauri.in/pankti/object"
)

func getStringFromArgs(arg object.Obj) (string, bool) { // Value, isOkay
	fileArg := arg

	if fileArg.Type() == object.STRING_OBJ {
		return fileArg.(*object.String).Value, true
	}
	return "", false

}

func ReadFile(args []object.Obj) object.Obj {
	filename, isOkay := getStringFromArgs(args[0])

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
	filename, isOkay := getStringFromArgs(args[0])

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
	filename, isOkay := getStringFromArgs(args[0])
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
	filename, isOkay := getStringFromArgs(args[0])

	if !isOkay {
		return &object.Error{Msg: "Filename must be string"}
	}

	if _, err := os.Stat(filename); err == nil {
		result = true
	}
	return &object.Boolean{Value: result}
}

func DeletePath(args []object.Obj) object.Obj {
	filename, ok := getStringFromArgs(args[0])

	if !ok {
		return &object.Error{Msg: "File path must be string"}
	}

	if _, err := os.Stat(filename); err != nil {
		return &object.Error{Msg: "File path does not exist"}
	} else {
		err := os.RemoveAll(filename)
		if err != nil {
			return &object.Error{Msg: "Failed to delete file"}
		}
	}
	return &object.Boolean{Value: true}

}

func RenameFile(args []object.Obj) object.Obj {

	if len(args) != 2 {
		return &object.Error{Msg: "Rename takes only two arguments"}
	}

	//result := false

	targetFile, isOkay := getStringFromArgs(args[0])

	if !isOkay {
		return &object.Error{Msg: "Target Filename must be string"}
	}

	newName, isOkay2 := getStringFromArgs(args[1])

	if !isOkay2 {
		return &object.Error{Msg: "New Filename must be string"}
	}

	if _, err := os.Stat(targetFile); err != nil {
		return &object.Error{Msg: "Target File does not exist"}
	}

	err := os.Rename(targetFile, newName)

	if err != nil {
		return &object.Error{Msg: "Failed to rename file"}
	}

	return &object.Boolean{Value: true}

}
