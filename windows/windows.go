//go:generate goversioninfo -64 -icon=res/icon.ico -manifest=res/pankti.exe.manifest
//go:generatei goversioninfo -icon=res/icon.ico -manifest=res/pankti.exe.manifest
package main

import (
	"os" 
	"go.cs.palashbauri.in/pankti/ide"
	log "github.com/sirupsen/logrus"
)
func init() {
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}

func main() {
	ide.RunIde()

}
