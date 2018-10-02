package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Prompt struct {
	*tview.InputField
	*History
}

func NewPrompt() *Prompt {
	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorDefault)

	p := &Prompt{
		InputField: inputField,
		History:    &History{},
	}

	return p
}

func (p *Prompt) SetDoneFunc(f func(tcell.Key)) {
	p.InputField.SetDoneFunc(func(key tcell.Key) {
		text := p.InputField.GetText()
		if len(text) > 0 {
			p.History.Add(text)
		}

		f(key)
	})
}

// InputHandler returns the handler for this primitive.
func (p *Prompt) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return p.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			hit, _ := p.SearchUp(p.GetText())
			p.SetText(hit)
		case tcell.KeyDown:
			hit, _ := p.SearchDown(p.GetText())
			p.SetText(hit)
		default:
			p.InputField.InputHandler()(event, setFocus)
		}
	})
}
