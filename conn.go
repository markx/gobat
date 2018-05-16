package main

import (
	"bufio"
	"net"

	"github.com/ziutek/telnet"
	"golang.org/x/net/proxy"
)

type Conn struct {
	net.Conn
	r *bufio.Reader
}

func Dial(addr string) (*Conn, error) {
	dialer := proxy.FromEnvironment()
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return NewConn(conn)
}

func NewConn(conn net.Conn) (*Conn, error) {
	telnetConn, err := telnet.NewConn(conn)
	if err != nil {
		return nil, err
	}

	c := Conn{
		Conn: telnetConn,
		r:    bufio.NewReaderSize(telnetConn, 1024),
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
