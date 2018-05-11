package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type SimpleUI struct {
	chatOutput io.WriteCloser
}

func NewSimpleUI(addr string) (*SimpleUI, error) {
	return &SimpleUI{}, nil
}

func (ui *SimpleUI) Run() error {
	inputs := make(chan string)

	chatFile, err := os.OpenFile("chat.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer chatFile.Close()
	ui.chatOutput = chatFile

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		inputs <- line
	}

	return nil
}

func (ui *SimpleUI) handleMessage(m Message) {
	if m.hasTag("chat") {
		log.Printf("msg chat: %#v", m)
		fmt.Fprint(ui.chatOutput, m.Content)
		return
	}

	fmt.Print(m.Content)
}
