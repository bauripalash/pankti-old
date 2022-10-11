//go:build !noide
// +build !noide

package ide

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/png"
	"os"
	"pankti/evaluator"
	"pankti/lexer"
	"pankti/object"
	"pankti/parser"
	"strings"

	"github.com/gen2brain/iup-go/iup"
)

//go:embed res/gear.png
var gearImg []byte

func SaveFile(filename string, content string, overwrite bool) error {
	if overwrite {
		err := os.WriteFile(filename, []byte(content), 0644)
		return err
	} else {
		nf, err := os.Create(filename)
		defer nf.Close()
		if err != nil {
			return err
		}

		nf.Write([]byte(content))
		nf.Sync()
	}

	return nil
}

func ShowAboutMenu() {
	iup.Message("About Pankti IDE", "Pankti IDE is a basic Editor to quickly run and test Pankti Programs")
	//iup.Popup(aboutDlg , iup.CENTER , iup.CENTER)

}

func ShowHelp() {
	iup.Message("Help with Pankti IDE", "TODO : Help")
}

func RunFile(src string) string {
	l := lexer.NewLexer(src)
	p := parser.NewParser(&l)

	prog := p.ParseProg()

	if len(p.GetErrors()) >= 1 {
		tempErrs := []string{}

		for _, item := range p.GetErrors() {
			tempErrs = append(tempErrs, item.String())
		}

		return strings.Join(tempErrs, " \n")
	}
	env := object.NewEnv()
	eh := evaluator.ErrorHelper{Source: src}
	evd := evaluator.Eval(prog, env, eh)

	if evd != nil {
		return evd.Inspect()
	} else {
		return ""
	}

}

func OpenFile(filename string) (string, error) {

	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func RunIde() {
	iup.Open()
	defer iup.Close()
	gearImage, err := png.Decode(bytes.NewReader(gearImg))
	if err != nil {
		fmt.Println("Failed to read gear image")
		return
	}

	iup.ImageFromImage(gearImage).SetHandle("gearimage")
	editor := iup.MultiLine().SetCallback("ACTION", iup.TextActionFunc(func(ih iup.Ihandle, item int, text string) int {
		if item == iup.K_g {
			return iup.IGNORE
		}

		return iup.DEFAULT
	})).SetAttributes(map[string]string{
		"EXPAND": "YES",
		"BORDER": "YES",
	})
	iup.SetGlobal("UTF8MODE", "YES")

	itemAbout := iup.Item("About").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		ShowAboutMenu()
		return iup.DEFAULT
	}))

	itemHelp := iup.Item("Help").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		ShowHelp()
		return iup.DEFAULT
	}))

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

	helpMenu := iup.Menu(itemHelp, itemAbout)
	helpSubmenu := iup.Submenu("Help", helpMenu)
	menu := iup.Menu(GetFileMenu(), helpSubmenu)
	menu.SetHandle("mymenu")

	iup.GetHandle("menuOpen").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
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

	iup.GetHandle("menuExit").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		return iup.CLOSE
	}))

	iup.GetHandle("menuSave").SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		iup.Popup(fopen, iup.CENTER, iup.CENTER)
		fopenStat := fopen.GetInt("STATUS")
		if fopenStat == 1 {
			err := SaveFile(fopen.GetAttribute("VALUE"), editor.GetAttribute("VALUE"), false)
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

	runBtn := iup.Button("/> Run ").SetAttribute("SIZE", "FIVExFIVE").SetAttribute("IMAGE", "gearimage")
	topToolbar := iup.Vbox(runBtn)
	outputBox := iup.MultiLine().SetAttributes(map[string]string{
		"EXPAND":   "YES",
		"BORDER":   "YES",
		"VALUE":    "output..",
		"READONLY": "YES",
	})

	runBtn.SetCallback("ACTION", iup.ActionFunc(func(i iup.Ihandle) int {
		outputBox.SetAttribute("VALUE", RunFile(editor.GetAttribute("VALUE")))
		return iup.IGNORE
	}))

	inf := iup.Vbox(
		topToolbar,
		editor,
		outputBox,
	)

	dlg := iup.Dialog(inf).SetAttributes(map[string]string{
		"MENU":  "mymenu",
		"TITLE": "Pankti IDE",
		"SIZE":  "QUARTERxQUARTER",
	})

	iup.Show(dlg)
	iup.MainLoop()
}
