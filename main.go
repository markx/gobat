package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ziutek/telnet"
)

func Com(t *telnet.Conn, c chan string) {
	for {
		rune, _, err := t.ReadRune()
		if err != nil {
			panic(err)
		}
		c <- string(rune)
	}
}

func Prompt(c chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		c <- line
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

	conn, err := telnet.Dial("tcp", dst)
	defer conn.Close()

	if err != nil {
		panic(err)
	}

	serverMessageChan := make(chan string)
	userInputChan := make(chan string)
	go Com(conn, serverMessageChan)
	go Prompt(userInputChan)

	for {
		select {
		case serverMessage := <-serverMessageChan:
			print(serverMessage)
		case userInput := <-userInputChan:
			conn.Write([]byte(userInput))
		}
	}
}
