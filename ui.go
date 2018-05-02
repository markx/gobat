package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type UI struct {
	client  *Client
	view    *tview.Application
	general *tview.TextView
	chat    *tview.TextView
	input   *tview.InputField
}

func NewUI(addr string) (*UI, error) {
	client := NewClient(addr)

	ui := &UI{
		client: client,
	}
	ui.initUI()

	return ui, nil
}

func (ui *UI) Run() error {
	errChan := make(chan error)

	go func() {
		if err := ui.client.Run(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		errChan <- ui.view.Run()
	}()

	for {
		select {
		case m := <-ui.client.Read():
			ui.handleMessage(m)
		case err := <-errChan:
			return err
		}
	}
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
		SetLabel(">").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorDefault)

	inputField.SetDoneFunc(func(key tcell.Key) {
		ui.client.Write(inputField.GetText() + "\n")
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatWindow, 15, 0, false).
		AddItem(generalWindow, 0, 1, false).
		AddItem(inputField, 1, 0, true)

	app.SetRoot(flex, true).SetFocus(flex)

	ui.view = app
	ui.input = inputField
	ui.general = generalWindow
	ui.chat = chatWindow
}

func (ui *UI) handleMessage(m Message) {
	if m.hasTag("chat") {
		fmt.Fprint(tview.ANSIIWriter(ui.chat), m.Content)
		return
	}
	fmt.Fprint(tview.ANSIIWriter(ui.general), m.Content)
}
