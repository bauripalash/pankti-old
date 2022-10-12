//go:build !noide
// +build !noide

package ide

import "github.com/gen2brain/iup-go/iup"
import _ "embed"

//go:embed res/ide_about.md
var AboutText string

func ShowAboutMenu() {
	aboutDlg := iup.Dialog(
		iup.Vbox(
			iup.Label(AboutText).SetAttributes(`EXPAND=YES,`),
			iup.Button("Quit").SetHandle("quitAbout").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
				return iup.CLOSE
			})).SetAttributes(`ALIGNMENT=ACENTER, MARGIN=10x10, GAP=10`),
		),
	).SetAttributes(`TITLE="About Pankti IDE", DEFAULTESC=quitAbout, SIZE=QUARTERxQUARTER`)
	iup.Show(aboutDlg)
	//iup.Message("About Pankti IDE", "Pankti IDE is a basic Editor to quickly run and test Pankti Programs")
    iup.Popup(aboutDlg , iup.CENTER , iup.CENTER)
	//iup.Popup(aboutDlg , iup.CENTER , iup.CENTER)

}

func ShowHelp() {
	iup.Message("Help with Pankti IDE", "TODO : Help")
}
