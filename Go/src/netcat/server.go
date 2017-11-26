package main

import (
	"bufio"
	"log"
	"net"
      "time"
      "strings"
      "fmt"

)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func echo(c net.Conn, shout string, delay time.Duration) {
      fmt.Fprintln(c , "\t" , strings.ToUpper(shout))
      time.Sleep(delay)
      fmt.Fprintln(c , "\t" , shout)
      time.Sleep(delay)
      fmt.Fprintln(c , "\t" , strings.ToLower(shout))
}

func handleConnection(c net.Conn) {
      input := bufio.NewScanner(c)
	for input.Scan() {
            echo(c, input.Text(), 1 * time.Second)
	}
      // NOTE: ignoring potential errors from input.Err()
}
