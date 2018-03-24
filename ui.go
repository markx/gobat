package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

var ui *UI

type UI struct {
	server *Server
	g      *gocui.Gui
	quit   chan struct{}
}

func NewUI(s *Server) *UI {
	ui = &UI{
		server: s,
		quit:   make(chan struct{}),
	}

	return ui
}

func (ui *UI) Run() error {
	err := ui.server.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to server: %v", err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		// handle error
		return fmt.Errorf("could not create ui: %v", err)
	}
	ui.g = g

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return fmt.Errorf("could not set key binding: %v", err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("error from main loop: %v", err)
	}

	if err := ui.Close(); err != nil {
		return fmt.Errorf("could not close gracefully: %v", err)
	}
	return nil
}

func (ui *UI) Close() error {
	close(ui.quit)
	ui.g.Close()
	return nil
}

func layout(g *gocui.Gui) error {
	if err := createMainView(g); err != nil {
		return fmt.Errorf("failed to create main view: %v", err)
	}

	if err := createInputView(g); err != nil {
		return fmt.Errorf("failed to create input view: %v ", err)
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
			for range time.Tick(100 * time.Millisecond) {
				select {
				case <-ui.quit:
					return
				default:
					updateMainView(g)
				}
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
	if len(msg) == 0 {
		return
	}

	g.Update(func(g *gocui.Gui) error {
		view, _ := g.View("mainView")
		_, err := view.Write(msg)
		if err != nil {
			return fmt.Errorf("could not write to main view: %v", err)
		}
		return nil

	})
}
