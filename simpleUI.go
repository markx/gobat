package main

import (
	"bufio"
	"fmt"
	"os"
)

type SimpleUI struct {
	client *Client
}

func NewSimpleUI(addr string) (*SimpleUI, error) {
	client, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	return &SimpleUI{
		client: client,
	}, nil
}

func (ui *SimpleUI) Run() error {
	cmds := make(chan string)

	go ui.client.Run()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			cmds <- line
		}
	}()

	for {
		select {
		case line := <-ui.client.Read():
			fmt.Print(line)
		case cmd := <-cmds:
			ui.client.Write(cmd)
		case err := <-ui.client.Errs():
			return err
		}
	}
}
