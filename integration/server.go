package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Str struct {
	St string
}

func main() {
	RunTCPServer()
}

func HandleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &Str{}
	dec.Decode(p)
	fmt.Println("gypsy")
	fmt.Println(p.St)
	fmt.Println("gypsy")
}

func RunTCPServer() {
	ln, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			log.Fatal(err)
		}
		go HandleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}
