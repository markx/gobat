package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var ui *UI

type UI struct {
	server *Server
	g      *gocui.Gui
}

func NewUI(s *Server) *UI {
	ui = &UI{
		server: s,
	}

	return ui
}

func (ui *UI) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		// handle error
		return fmt.Errorf("could not create ui: %v", err)
	}
	defer g.Close()

	ui.g = g

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return fmt.Errorf("error from main loop: %v", err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("error from main loop: %v", err)
	}

	return nil
}

func layout(g *gocui.Gui) error {
	if err := createMainView(g); err != nil {
		log.Print("failed to create layout", err)
		return err
	}

	if err := createInputView(g); err != nil {
		log.Print("failed to create layout", err)
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func createMainView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("mainView", 0, 0, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Autoscroll = true
		v.Wrap = true

		go func() {
			for {
				updateMainView(g)
			}
		}()

	}

	return nil
}

func createInputView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("inputView", 0, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err := g.SetCurrentView("inputView")
		if err != nil {
			return err
		}

		v.Frame = true
		v.Editable = true
		v.Autoscroll = false
		v.Editor = gocui.EditorFunc(simpleEditor)
	}

	return nil
}

func updateMainView(g *gocui.Gui) {
	msg, err := ui.server.Read()
	if err != nil {
		log.Println(err)
		return
	}

	g.Execute(func(g *gocui.Gui) error {
		view, _ := g.View("mainView")
		fmt.Fprint(view, string(msg))
		return nil
	})
}
