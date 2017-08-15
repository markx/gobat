package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	host := "aardmud.org"
	port := "4000"

	if len(args) >= 2 {
		host = args[0]
		port = args[1]
	}

	server := &Server{
		host: host,
		port: port,
	}

	err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	gui := NewUI(server)
	gui.Run()

}
