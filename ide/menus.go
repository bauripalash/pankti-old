//go:build !noide
// +build !noide

package ide

import "github.com/gen2brain/iup-go/iup"

func GetFileMenu() iup.Ihandle {

	itemOpen := iup.Item("Open").SetAttribute("KEY", "O").SetHandle("menuOpen")
	itemSave := iup.Item("Save").SetAttribute("KEY", "S").SetHandle("menuSave")
	itemUndo := iup.Item("Undo").SetAttributes("KEY=U, ACTIVE=NO").SetHandle("menuUndo")
	itemExit := iup.Item("Exit").SetAttribute("KEY", "x").SetHandle("menuExit")

	fileMenu := iup.Menu(itemOpen, itemSave, iup.Separator(), itemUndo, itemExit)
	subMenu := iup.Submenu("File", fileMenu)

	return subMenu
}

func GetHelpMenu() iup.Ihandle {
	itemAbout := iup.Item("About").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		ShowAboutMenu()
		return iup.DEFAULT
	}))

	itemHelp := iup.Item("Help").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		ShowHelp()
		return iup.DEFAULT
	}))

	helpMenu := iup.Menu(itemHelp, itemAbout)
	helpSubmenu := iup.Submenu("Help", helpMenu)

	return helpSubmenu
}
