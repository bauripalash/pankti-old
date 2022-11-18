//go:build !noide
// +build !noide

package ide

import "github.com/gen2brain/iup-go/iup"
import _ "embed"

//go:embed res/ide_about.md
var AboutText string

//go:embed res/ide_help.md
var HelpText string

func ShowAboutMenu() {

	iup.Message("About Pankti IDE", AboutText)

}

func ShowHelp() {
	iup.Message("Help with Pankti IDE", HelpText)
}
