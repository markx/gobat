package main

import (
	"fmt"
)

type Client struct {
	addr     string
	conn     *Conn
	messages chan string
	cmds     chan string
	errs     chan error
}

func NewClient(addr string) (*Client, error) {
	c := &Client{
		addr:     addr,
		messages: make(chan string),
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

	quit := make(chan struct{})
	errChan := make(chan error)
	defer close(quit)

	go func() {
		for {
			select {
			case cmd := <-c.cmds:
				_, err := c.conn.Write([]byte(cmd))
				if err != nil {
					errChan <- err
					return
				}
			case <-quit:
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
				return fmt.Errorf("Failed to read:%v", err)
			}
			c.messages <- line
		}
	}
}

func (c *Client) Write(cmd string) {
	c.cmds <- string(cmd)
}

func (c *Client) Read() <-chan string {
	return c.messages
}

func (c *Client) Errs() <-chan error {
	return c.errs
}
