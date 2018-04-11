package main

import (
	"fmt"
)

type Client struct {
	addr     string
	conn     *Conn
	messages chan Message
	cmds     chan string
}

func NewClient(addr string) (*Client, error) {
	c := &Client{
		addr:     addr,
		messages: make(chan Message),
		cmds:     make(chan string),
	}

	return c, nil
}

func (c *Client) Run() error {
	conn, err := Dial(c.addr)
	if err != nil {
		return err
	}
	c.conn = conn

	errChan := make(chan error)

	go func() {
		for cmd := range c.cmds {
			_, err := c.conn.Write([]byte(cmd))
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	for {
		select {
		case err := <-errChan:
			return err
		default:
			line, err := c.conn.ReadLine()
			if err != nil {
				return fmt.Errorf("Failed to read: %v", err)
			}
			m := NewMessage(line)
			c.messages <- m
		}
	}
}

func (c *Client) Write(cmd string) {
	c.cmds <- string(cmd)
}

func (c *Client) Read() <-chan Message {
	return c.messages
}
