package stdlib

import (
	"io/fs"
	"os"
	"path/filepath"

	"go.cs.palashbauri.in/pankti/object"
)

func getStringFromArgs(arg object.Obj) (string, bool) { // Value, isOkay
	fileArg := arg

	if fileArg.Type() == object.STRING_OBJ {
		return fileArg.(*object.String).Value, true
	}
	return "", false

}

func ReadFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
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

func CreateEmptyFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
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

func WriteToFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
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

func FileDirExists(eh *object.ErrorHelper, args []object.Obj) object.Obj {
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

func DeletePath(eh *object.ErrorHelper, args []object.Obj) object.Obj {
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

func RenameFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {

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

func IsAFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	target, ok := getStringFromArgs(args[0])
	result := false

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, "Target Filename must be string")
	}

	if s, err := os.Stat(target); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, "Target does not exist")
	} else {
		if !s.IsDir() {
			result = true
		}
	}

	return &object.Boolean{Value: result}
}

func IsADir(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	target, ok := getStringFromArgs(args[0])
	result := false

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, "Target Filename must be string")
	}

	if s, err := os.Stat(target); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, "Target does not exist")
	} else {
		if s.IsDir() {
			result = true
		}
	}

	return &object.Boolean{Value: result}
}

func AppendLineToFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	filename, ok := getStringFromArgs(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, "Target Filename must be string")
	}

	data, ok2 := getStringFromArgs(args[1])

	if !ok2 {
		return object.NewErr(args[1].GetToken(), eh, true, "Data must be a string")
	}

	if s, err := os.Stat(filename); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, "Target does not exist")
	} else {
		if s.IsDir() {
			return object.NewErr(args[0].GetToken(), eh, true, "Target is a directory")
		} else {
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)

			if err != nil {

				return object.NewErr(args[0].GetToken(), eh, true, "Failed to open file")
			}

			defer f.Close()

			if _, err := f.WriteString(data); err != nil {
				return object.NewErr(args[0].GetToken(), eh, true, "Failed to write data to file ")
			}
		}

	}

	return &object.Boolean{Value: true}

}

func ListDir(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	d := args[0]
	dirname, ok := getStringFromArgs(d)
	if !ok {
		return object.NewErr(d.GetToken(), eh, true, "Target must be string")
	}

	if f, err := os.Stat(dirname); err != nil {
		return object.NewErr(d.GetToken(), eh, true, "Target does not exist")
	} else {
		if !f.IsDir() {
			return object.NewErr(d.GetToken(), eh, true, "Target is not a directory")
		}

		result := []object.Obj{}

		err := filepath.Walk(dirname, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			result = append(result, &object.String{Value: path})
			return nil
		})

		if err != nil {
			object.NewErr(d.GetToken(), eh, true, "Failed to list files in directory")
		}

		return &object.Array{Elms: result}

	}

}
