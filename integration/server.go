package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	port, _ := strconv.Atoi(os.Args[1])
	laddr := &net.TCPAddr{
		IP:   net.IP("127.0.0.1"),
		Port: port,
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal("listen tcp failed")
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal("tcp accept failed", err)
		}
		b := []byte{}
		_, _ = conn.Read(b)
		fmt.Println(b)
	}
}
