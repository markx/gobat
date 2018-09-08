package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type UI struct {
	tviewApp     *tview.Application
	general      *tview.TextView
	chat         *tview.TextView
	input        *tview.InputField
	inputHandler func(string)
}

func NewUI() *UI {
	ui := &UI{}
	ui.initUI()

	return ui
}

func (ui *UI) Run() error {
	return ui.tviewApp.Run()
}

func (ui *UI) Stop() {
	ui.tviewApp.Stop()
}

func (ui *UI) initUI() {
	app := tview.NewApplication()

	generalWindow := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorDefault).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	generalWindow.SetBorder(true)

	chatWindow := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorDefault).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	chatWindow.SetBorder(true)

	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorDefault)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if ui.inputHandler == nil {
			return
		}
		ui.inputHandler(inputField.GetText())
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatWindow, 15, 0, false).
		AddItem(generalWindow, 0, 1, false).
		AddItem(inputField, 1, 0, true)

	app.SetRoot(flex, true).SetFocus(flex)

	ui.tviewApp = app
	ui.input = inputField
	ui.general = generalWindow
	ui.chat = chatWindow
}

func (ui *UI) SendToWindow(window, content string) {
	switch window {
	case "chat":
		fmt.Fprint(tview.ANSIIWriter(ui.chat), content)
	default:
		fmt.Fprint(tview.ANSIIWriter(ui.general), content)
	}
}

func (ui *UI) SetInputHandler(handler func(string)) {
	ui.inputHandler = handler
}
