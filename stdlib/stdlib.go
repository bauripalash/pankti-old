package stdlib

import (
	"embed"
	"errors"
	"go.cs.palashbauri.in/pankti/constants"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed x
var stdx embed.FS

func GetStdLibFileSrc(path string) (string, bool) {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
			f, er := os.ReadFile(path)
			if er == nil {
				return string(f), true
			}

		}
	} else {
		cwd, _ := os.Getwd()
		p := filepath.Join(cwd, path)
		if _, err := os.Stat(p); !errors.Is(err, fs.ErrNotExist) {
			f, er := os.ReadFile(p)
			if er == nil {
				return string(f), true
			}
		}
		enName, ok := constants.GetStdName(path)
		if ok {
			stdpath := "x/" + enName
			f, err := stdx.ReadFile(stdpath)
			//fmt.Println(err.Error())
			if err == nil {
				return string(f), true
			}
		}
	}

	return "", false
}
