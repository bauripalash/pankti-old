//go:generate goversioninfo -64 -icon=windows/res/icon.ico -manifest=windows/res/pankti.exe.manifest
package main

import (
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"

	"go.cs.palashbauri.in/pankti/cmd"
	"go.cs.palashbauri.in/pankti/constants"

	log "github.com/sirupsen/logrus"
)

const DEBUG = true
const IMPORT_PATH_ENV = constants.IMPORT_PATH_ENV
const WINDOWS_IMPORT_PATH = "??" //TODO: Test on Windows
const LINUX_IMPORT_PATH = ".pankti/stdlib/"

func init() {
	//log.SetLevel(log.DebugLevel)
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	setImportPathEnv()
}

func setImportPathEnv() {
	custom_import_path := os.Getenv(IMPORT_PATH_ENV)

	if len(custom_import_path) < 1 && !DEBUG {

		if runtime.GOOS == "linux" {
			p, err := os.UserHomeDir()
			if err != nil {
				return
			}

			os.Setenv(IMPORT_PATH_ENV, path.Join(p, LINUX_IMPORT_PATH))

		} else if runtime.GOOS == "windows" {

			p, err := os.UserHomeDir()
			if err != nil {
				return
			}
			os.Setenv(IMPORT_PATH_ENV, path.Join(p, WINDOWS_IMPORT_PATH))

		}

	} else {
		curdir, _ := os.Getwd()
		os.Setenv(IMPORT_PATH_ENV, path.Join(curdir, "/stdlib/x/"))
	}
}

func main() {
	is_noide := false
	bi, noerr := debug.ReadBuildInfo()
	if !noerr {
		return
	}
	for _, item := range bi.Settings {
		if item.Key == "-tags" && strings.Contains(item.Value, "noide") {
			is_noide = true
			break
		}
	}

	cmd.Execute(is_noide)
}
