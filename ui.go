package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/markx/gobat/prompt"
	"github.com/rivo/tview"
)

type UI struct {
	tviewApp     *tview.Application
	main         *tview.TextView
	chat         *tview.TextView
	character    *tview.TextView
	input        *prompt.Prompt
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

	mainWindow := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorDefault).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	mainWindow.SetBorder(true).
		SetTitle("Main")

	chatWindow := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorDefault).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	chatWindow.SetBorder(true).
		SetTitle("Chat")

	characterWindow := tview.NewTextView()
	characterWindow.SetBorder(true).SetTitle("Player")

	input := prompt.NewPrompt()
	input.SetDoneFunc(func(key tcell.Key) {
		if ui.inputHandler == nil {
			return
		}
		ui.inputHandler(input.GetText())
	})

	sideWindow := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatWindow, 0, 1, false).
		AddItem(characterWindow, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(mainWindow, 128, 0, false).
			AddItem(sideWindow, 0, 1, false), 0, 1, false).
		AddItem(input, 1, 0, true)

	app.SetRoot(flex, true).SetFocus(flex)

	ui.tviewApp = app
	ui.input = input
	ui.main = mainWindow
	ui.chat = chatWindow
	ui.character = characterWindow
}

func (ui *UI) SendToWindow(window, content string) {
	switch window {
	case "chat":
		fmt.Fprint(tview.ANSIIWriter(ui.chat), content)
	default:
		fmt.Fprint(tview.ANSIIWriter(ui.main), content)
	}
}

func (ui *UI) SetInputHandler(handler func(string)) {
	ui.inputHandler = handler
}
