package main

import (
	"flag"
	"fmt"
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

	host := flag.String("host", "bat.org", "mud server host")
	port := flag.String("port", "23", "mud server port")
	flag.Parse()

	client, err := NewClient(fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Run(); err != nil {
		log.Fatal(err)
	}
}
