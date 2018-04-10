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
	inputs := make(chan string)
	errChan := make(chan error)

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
		case line := <-ui.client.Read():
			fmt.Print(line)
		case input := <-inputs:
			ui.client.Write(input)
		case err := <-errChan:
			return err
		}
	}
}
