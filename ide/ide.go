//go:build !noide
// +build !noide

package ide

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/png"

	"github.com/gen2brain/iup-go/iup"
)

//go:embed res/gear.png
var gearImg []byte

//go:embed res/icon.png
var iconImg []byte

func RunIde() {
	iup.Open()
	defer iup.Close()
	gearImage, err := png.Decode(bytes.NewReader(gearImg))
	if err != nil {
		fmt.Println("Failed to read gear image")
		return
	}
	iconImage, err := png.Decode(bytes.NewReader(iconImg))

	if err != nil {
		fmt.Println("Failed to read Icon Image")
		return
	}

	iup.ImageFromImage(gearImage).SetHandle("gearimage")
	iup.ImageFromImage(iconImage).SetHandle("iconimage")
	editor := iup.MultiLine().
		SetCallback("ACTION", iup.TextActionFunc(func(ih iup.Ihandle, item int, text string) int {
			return iup.DEFAULT
		})).
		SetAttributes(map[string]string{
			"EXPAND": "YES",
			"BORDER": "YES",
		})
	iup.SetGlobal("UTF8MODE", "YES")

	fd := iup.FileDlg().SetAttributes(map[string]string{
		"TITLE": "Open File",
	})
	defer fd.Destroy()

	fopen := iup.FileDlg().SetAttributes(map[string]string{
		"TITLE":      "Save File",
		"DIALOGTYPE": "SAVE",
		"FILTER":     "*.pank",
		"FILTERINFO": "Pankti Script",
	})

	defer fopen.Destroy()

	menu := iup.Menu(GetFileMenu(), GetHelpMenu())
	menu.SetHandle("mymenu")

	iup.GetHandle("menuOpen").
		SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
			iup.Popup(fd, iup.CENTER, iup.CENTER)

			if fd.GetInt("STATUS") == 0 {
				f, err := OpenFile(fd.GetAttribute("VALUE"))
				if err != nil {
					iup.Message("Error!", err.Error())
					return iup.IGNORE
				}
				editor.SetAttribute("VALUE", f)

			} else {
				fmt.Println(fd.GetInt("STATUS"))
			}

			return iup.IGNORE
		}))

	iup.GetHandle("menuExit").
		SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
			return iup.CLOSE
		}))

	iup.GetHandle("menuSave").
		SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
			iup.Popup(fopen, iup.CENTER, iup.CENTER)
			fopenStat := fopen.GetInt("STATUS")
			if fopenStat == 1 {
				err := SaveFile(
					fopen.GetAttribute("VALUE"),
					editor.GetAttribute("VALUE"),
					false,
				)
				if err != nil {
					fmt.Println("Failed to save file")

				}
			} else if fopenStat == 0 {
				switch iup.Alarm("Overwrite File?", "overwrite file "+editor.GetAttribute("VALUE")+" ?", "Okay", "No", "Cancel") {
				case 1:
					SaveFile(fopen.GetAttribute("VALUE"), editor.GetAttribute("VALUE"), true)
				default:
					iup.Message("File not Saved", "File not Overwritten")
				}
				//iup.Message("Save File" , "File was overwriten")
			}
			return iup.IGNORE
		}))

	runBtn := iup.Button("/> Run ").
		SetAttribute("SIZE", "FIVExFIVE").
		SetAttribute("IMAGE", "gearimage")
	topToolbar := iup.Vbox(runBtn)
	outputBox := iup.MultiLine().SetAttributes(map[string]string{
		"EXPAND":   "YES",
		"BORDER":   "YES",
		"VALUE":    "output..",
		"READONLY": "YES",
	})

	runBtn.SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		outputBox.SetAttribute(
			"VALUE",
			RunFile(editor.GetAttribute("VALUE")),
		)
		return iup.IGNORE
	}))

	splitbox := iup.Split(editor, outputBox).SetAttributes(map[string]string{
		"ORIENTATION": "HORIZONTAL",
	})

	inf := iup.Vbox(
		topToolbar,
		splitbox,
		/*
			editor,
			outputBox,
		*/
	)

	dlg := iup.Dialog(inf).SetAttributes(map[string]string{
		"MENU":  "mymenu",
		"TITLE": "Pankti IDE",
		"SIZE":  "QUARTERxQUARTER",
		"ICON":  "iconimage",
	})
	iup.SetCallback(dlg, "K_ANY", iup.KAnyFunc(func(ih iup.Ihandle, id int) int {
		if id == iup.K_F5 {
			outputBox.SetAttribute(
				"VALUE",
				RunFile(editor.GetAttribute("VALUE")))
			return iup.IGNORE

		}
		return iup.CONTINUE
	}))

	iup.Show(dlg)
	iup.MainLoop()
}
