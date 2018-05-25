package main

import (
	"fmt"
	"log"
)

type Client struct {
	conn     *Conn
	ui       *UI
	triggers *Triggers
}

func NewClient(addr string) (*Client, error) {
	conn, err := Dial(addr)
	if err != nil {
		return nil, err
	}

	ui := NewUI()

	c := &Client{
		conn:     conn,
		ui:       ui,
		triggers: NewTriggers(),
	}

	return c, nil
}

func (c *Client) handleInput(content string) {
	c.conn.Write([]byte(content + "\n"))
}

func (c *Client) Run() error {
	c.ui.SetInputHandler(c.handleInput)

	errChan := make(chan error, 1)

	go func() {
		errChan <- c.ui.Run()
	}()

	for {
		select {
		case err := <-errChan:
			return err
		default:
			line, err := c.conn.ReadLine()
			if err != nil {
				c.ui.Stop()
				return err
			}
			log.Printf("msg: %s", line)
			c.handleMessage(NewMessage(line))
		}
	}
}

func (c *Client) Send(cmd string) {
	fmt.Fprint(c.conn, cmd)
}

func (c *Client) handleMessage(m Message) {
	c.triggers.Match(&m, c)

	if m.hasTag("chat") {
		c.ui.SendToWindow("chat", m.Content)
		return
	}
	if m.hasTag("prompt") {
		c.ui.SendToWindow("general", m.Content)
		return
	}
	c.ui.SendToWindow("general", m.Content)
}
