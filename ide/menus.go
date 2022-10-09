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
