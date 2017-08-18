package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

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
