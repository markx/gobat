package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type SimpleUI struct {
	client     *Client
	chatOutput io.WriteCloser
}

func NewSimpleUI(addr string) (*SimpleUI, error) {
	client := NewClient(addr)

	return &SimpleUI{
		client: client,
	}, nil
}

func (ui *SimpleUI) Run() error {
	inputs := make(chan string)
	errChan := make(chan error)

	chatFile, err := os.OpenFile("chat.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer chatFile.Close()
	ui.chatOutput = chatFile

	go func() {
		if err := ui.client.Run(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			inputs <- line
		}
	}()

	for {
		select {
		case m := <-ui.client.Read():
			ui.handleMessage(m)
		case input := <-inputs:
			ui.client.Write(input)
		case err := <-errChan:
			return err
		}
	}
}

func (ui *SimpleUI) handleMessage(m Message) {
	if m.hasTag("chat") {
		log.Printf("msg chat: %#v", m)
		fmt.Fprint(ui.chatOutput, m.Content)
		return
	}

	fmt.Print(m.Content)
}
