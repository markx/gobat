package main

import (
	"fmt"
	"net"
	"time"

	"github.com/ziutek/telnet"
)

type Client struct {
	addr     string
	conn     net.Conn
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
	conn, err := telnet.DialTimeout("tcp", c.addr, 5*time.Second)
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

	buf := make([]byte, 1024)
	for {
		select {
		case err := <-errChan:
			return err
		default:
			n, err := c.conn.Read(buf)
			if err != nil {
				return fmt.Errorf("Failed to read:%v", err)
			}
			c.messages <- string(buf[:n])
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
