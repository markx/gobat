package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/jroimartin/gocui"
)

func simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyEnter:
		line := v.ViewBuffer()
		v.Clear()
		v.SetCursor(0, 0)

		// remove the extra space
		//this is a bug in gocui: https://github.com/jroimartin/gocui/issues/69
		line = strings.TrimRightFunc(line, unicode.IsSpace) + "\n"

		err := ui.server.Write([]byte(line))
		if err != nil {
			log.Panicln(err)
		}

		ui.g.Execute(func(g *gocui.Gui) error {
			view, _ := g.View("mainView")
			_, err := view.Write([]byte(line))
			if err != nil {
				return fmt.Errorf("could not write to main view: %v", err)
			}
			return nil
		})

	case key == gocui.KeyTab:
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyArrowDown:
	case key == gocui.KeyArrowUp:
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		cx, _ := v.Cursor()
		line := v.ViewBuffer()
		// if cx == 0 {
		// v.MoveCursor(-1, 0, false)
		if cx < len(line)-1 {
			v.MoveCursor(1, 0, false)
		}

	case key == gocui.KeyCtrlA:
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
	case key == gocui.KeyCtrlE:
		v.SetCursor(len(v.Buffer())-1, 0)
	case key == gocui.KeyCtrlK:
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
	case key == gocui.KeyCtrlLsqBracket:
		// logger.Logger.Println("word...")
	}

}
