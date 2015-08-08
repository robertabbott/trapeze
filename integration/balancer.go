package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/robertabbott/trapeze"
)

type Str struct {
	St string
}

type LoadBalancer struct {
	Addr *net.TCPAddr
}

func main() {
	lb := LoadBalancer{
		Addr: &net.TCPAddr{
			IP:   []byte("127.0.0.1"),
			Port: 3069,
		},
	}
	endpoints := []trapeze.ServiceEndpoint{}
	endpoints, _ = lb.AddEndpoint("127.0.0.1", 6969, endpoints)
	endpoints, _ = lb.AddEndpoint("127.0.0.1", 6868, endpoints)

	p := trapeze.Proxy{
		Addr:         lb.Addr,
		Endpoints:    endpoints,
		LoadBalancer: lb,
		Connections:  make(map[trapeze.ServiceEndpoint]map[*net.Addr]trapeze.Connection),
	}
	fmt.Println(p.Endpoints)

	fmt.Println("starting servers")
	// start servers on 6969 and 6868
	go RunTCPServer(6969)
	go RunTCPServer(6868)

	// send messages to each server for no reason
	SendStructTCP("127.0.0.1:6969", Str{"seamus"})
	SendStructTCP("127.0.0.1:6868", Str{"seamus"})
	time.Sleep(1 * time.Second)

	fmt.Println("starting loadbalancer")
	// run load balancer to route traffic to both servers
	go p.Layer4()

	for {
		// load balancer should route requests to 6868 and 6969
		SendStructTCP("127.0.0.1:3069", Str{"seamus"})
		time.Sleep(500 * time.Millisecond)
	}
}

func RunTCPServer(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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

func HandleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &Str{}
	dec.Decode(p)
	fmt.Println(conn.RemoteAddr())
	fmt.Println(p.St)
	fmt.Println(conn.LocalAddr())
}

// pick endpoint from list at random
func (l LoadBalancer) NextEndpoint(endpoints []trapeze.ServiceEndpoint, client trapeze.Connection) trapeze.Connection {
	index := rand.Intn(len(endpoints))
	return populateConn(&endpoints[index], client)
}

func (l LoadBalancer) AddEndpoint(addr string, port int, endpoints []trapeze.ServiceEndpoint) ([]trapeze.ServiceEndpoint, error) {
	endpoints = append(endpoints, trapeze.ServiceEndpoint{
		Addr: &net.TCPAddr{
			IP:   []byte(addr),
			Port: port,
		},
		Port: port,
	})
	return endpoints, nil
}

func (l LoadBalancer) RemoveEndpoint(endpoint trapeze.ServiceEndpoint) {
	return
}

func populateConn(endpoint *trapeze.ServiceEndpoint, client trapeze.Connection) trapeze.Connection {
	client.RouteTo = endpoint
	client.CloseCh = make(chan struct{})
	return client
}

func SendStructTCP(addr string, s Str) {
	conn := ConnectTCP(addr)
	if conn == nil {
		return
	}
	sendStruct(&s, conn)
	conn.Close()
}

func ConnectTCP(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil
	}
	return conn
}

func sendStruct(st *Str, conn net.Conn) {
	gob.Register(st.St)
	enc := gob.NewEncoder(conn)
	err := enc.Encode(st)
	if err != nil {
		log.Fatal(err)
	}
}

func ShutdownServer(conn *net.TCPListener) {
	conn.Close()
}
