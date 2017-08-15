package main

import (
	"fmt"
	"time"

	"github.com/tsuru/tsuru/log"
	"github.com/ziutek/telnet"
)

type Server struct {
	host string
	port string

	connection *telnet.Conn
}

func (s *Server) Connect() error {
	dst := s.host + ":" + s.port
	conn, err := telnet.DialTimeout("tcp", dst, 5*time.Second)
	if err != nil {
		return fmt.Errorf("could not dial: %v", err)
	}

	s.connection = conn
	return nil
}

func (s *Server) Close() {
	err := s.connection.Close()
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) Write(content []byte) error {
	_, err := s.connection.Write(content)
	return err
}

func (s *Server) Read() ([]byte, error) {
	byteChan := make(chan byte)
	errChan := make(chan error)

	quit := make(chan struct{})
	defer close(quit)

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				b, err := s.connection.ReadByte()
				if err != nil {
					errChan <- err
					return
				}
				byteChan <- b
			}
		}
	}()

	var line []byte
	for {
		select {
		case b := <-byteChan:
			line = append(line, b)
			if b == '\n' {
				return line, nil
			}
		case err := <-errChan:
			return line, err

		case <-time.After(time.Millisecond * 300):
			return line, nil
		}
	}
}
