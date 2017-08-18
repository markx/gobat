package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	host := flag.String("host", "aardmud.org", "mud server host")
	port := flag.String("port", "4000", "mud server port")
	flag.Parse()

	server := &Server{
		host: *host,
		port: *port,
	}

		log.Fatal(err)
	}

	gui := NewUI(server)
	gui.Run()

}
