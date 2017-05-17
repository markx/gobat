package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ziutek/telnet"
)

func Com(t *telnet.Conn) {
	for {
		buf, err := t.ReadString('\n')
		if err != nil {
			panic(err)
		}

		print(buf)
	}

}

func prompt(t io.Writer) {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		t.Write(line)
	}
}

func print(s ...interface{}) {
	fmt.Print(s...)
}

func main() {
	args := os.Args[1:]

	host := "aardmud.org"
	port := "4000"

	if len(args) >= 2 {
		host = args[0]
		port = args[1]
	}

	dst := host + ":" + port

	t, err := telnet.Dial("tcp", dst)
	defer t.Close()

	if err != nil {
		panic(err)
	}

	go Com(t)
	prompt(t)
}
