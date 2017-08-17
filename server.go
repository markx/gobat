package main

import (
	"fmt"
	"log"
	"time"

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
	conn.SetUnixWriteMode(true)

	s.connection = conn
	return nil
}

func (s *Server) Close() error {
	err := s.connection.Close()
	if err != nil {
		return fmt.Errorf("could not close: %v", err)
	}
	return nil
}

func (s *Server) Write(content []byte) error {
	_, err := s.connection.Write(content)
	return err
}

func (s *Server) Read() ([]byte, error) {
	buf := make([]byte, 1024)
	b, err := s.connection.Read(buf)
	log.Println(b, err)
	return buf[:b], err
}
