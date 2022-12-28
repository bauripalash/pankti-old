package stdlib

import (
	"io/fs"
	"os"
	"path/filepath"

	"go.cs.palashbauri.in/pankti/errs"
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
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	filename, isOkay := getStringFromArgs(args[0])

	if !isOkay {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}
	d, err := os.ReadFile(filename)
	if err != nil {
		return &object.Error{Msg: "Failed to read file"}
	}

	return &object.String{Value: string(d)}
}

func CreateEmptyFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	filename, isOkay := getStringFromArgs(args[0])

	if !isOkay {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	if _, err := os.Stat(filename); err == nil {
		return &object.Boolean{Value: false}
	}

	f, err := os.Create(filename)
	if err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FAILED_TO_CREATE"], filename)

	}

	if err := f.Close(); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FAILED_TO_CLOSE_FILE"], filename)
	} else {
		return &object.Boolean{Value: true}
	}

}

func WriteToFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	filename, isOkay := getStringFromArgs(args[0])
	data := args[1]
	if !isOkay {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	err := os.WriteFile(filename, []byte(data.Inspect()), 0644)

	if err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FAILED_TO_WRITE_FILE"], filename)
	}

	return &object.Boolean{Value: true}

}

func FileDirExists(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	result := false
	filename, isOkay := getStringFromArgs(args[0])

	if !isOkay {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	if _, err := os.Stat(filename); err == nil {
		result = true
	}
	return &object.Boolean{Value: result}
}

func DeletePath(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	filename, ok := getStringFromArgs(args[0])

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	if _, err := os.Stat(filename); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_NOT_EXIST"])
	} else {
		err := os.RemoveAll(filename)
		if err != nil {
			return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["DELETE_FAILED"])
		}
	}
	return &object.Boolean{Value: true}

}

func RenameFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	//result := false

	targetFile, isOkay := getStringFromArgs(args[0])

	if !isOkay {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	newName, isOkay2 := getStringFromArgs(args[1])

	if !isOkay2 {
		return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	if _, err := os.Stat(targetFile); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_NOT_EXIST"])
	}

	err := os.Rename(targetFile, newName)

	if err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["RENAME_FAILED"])
	}

	return &object.Boolean{Value: true}

}

func IsAFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	target, ok := getStringFromArgs(args[0])
	result := false

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	if s, err := os.Stat(target); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_NOT_EXIST"])
	} else if !s.IsDir() {
		result = true
	}

	return &object.Boolean{Value: result}
}

func IsADir(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	target, ok := getStringFromArgs(args[0])
	result := false

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	if s, err := os.Stat(target); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_NOT_EXIST"])
	} else if s.IsDir() {
		result = true
	}

	return &object.Boolean{Value: result}
}

func AppendLineToFile(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	filename, ok := getStringFromArgs(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_PATH_MUST_BE_STRING"])
	}

	data, ok2 := getStringFromArgs(args[1])

	if !ok2 {
		return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["DATA_MUST_BE_STRING"])
	}

	if s, err := os.Stat(filename); err != nil {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FILE_NOT_EXIST"])
	} else {
		if s.IsDir() {
			return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TARGET_IS_DIR"])
		} else {
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)

			if err != nil {

				return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FAILED_OPEN_FILE"])
			}

			defer f.Close()

			if _, err := f.WriteString(data); err != nil {
				return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["FAILED_TO_WRITE_DATA"])
			}
		}

	}

	return &object.Boolean{Value: true}

}

func ListDir(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if IsAndroid() {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ON_ANDROID"])
	}
	d := args[0]
	dirname, ok := getStringFromArgs(d)
	if !ok {
		return object.NewErr(d.GetToken(), eh, true, errs.Errs["FILENAME_MUST_BE_STRING"])
	}

	if f, err := os.Stat(dirname); err != nil {
		return object.NewErr(d.GetToken(), eh, true, errs.Errs["FILE_NOT_EXIST"])
	} else {
		if !f.IsDir() {
			return object.NewErr(d.GetToken(), eh, true, errs.Errs["TARGET_NO_DIR"])
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
			object.NewErr(d.GetToken(), eh, true, errs.Errs["DIR_LIST_FAILED"])
		}

		return &object.Array{Elms: result}

	}

}
