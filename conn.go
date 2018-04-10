package main

import (
	"bufio"
	"net"
	"time"

	"github.com/ziutek/telnet"
)

type Conn struct {
	net.Conn
	r *bufio.Reader
}

func Dial(addr string) (*Conn, error) {
	conn, err := telnet.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

func NewConn(conn net.Conn) (*Conn, error) {
	c := Conn{
		Conn: conn,
		r:    bufio.NewReaderSize(conn, 1024),
	}
	return &c, nil
}

func (c *Conn) ReadLine() (string, error) {
	var line []byte

	for {
		b, err := c.r.ReadByte()
		if err != nil {
			return string(line), err
		}

		line = append(line, b)

		if b == '\n' || c.r.Buffered() == 0 {
			return string(line), nil
		}
	}
}
