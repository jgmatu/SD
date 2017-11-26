package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.45:4556")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "Hi world!\n")
	defer conn.Close()
}
