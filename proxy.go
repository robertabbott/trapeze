package trapeze

import (
	"fmt"
	"io"
	"net"
)

// implements OSI layer 3 load balancer
// layer three load balancing occurs at the dns
// level. If you app/site/whatever has multiple
// external IPs a layer 3 load balancer will route
// traffic between the 3 IPs without inspecting packets
// for application specific information

// reverse proxy takes internet traffic and routes in to
// members of the routing pool. The 'routing pool' in this
// case is a slice of ServiceEndpoints

type Proxy struct {
	Addr         *net.TCPAddr
	Endpoints    []ServiceEndpoint
	loadBalancer LoadBalancer
	connections  map[ServiceEndpoint][]connection
}

// keep slice of connections

func (p *Proxy) Layer3() {
	route := make(chan connection)

	// listen for incoming connections
	// schedule ServiceEndpoints based on some algorithm

	go listenForConnections(p.Addr, route)
	for {
		select {
		case c := <-route:
			conn := p.loadBalancer.NextEndpoint(&p.Endpoints, c)
			addConn(p.connections, conn)
			go routeConn(conn, ch)
		}
	}
}

func addConn(connections *map[ServiceEndpoint][]connection, conn connection) {
	connections[conn.routeTo] = append(connections[conn.routeTo], conn)
}

func routeConn(c connection, ch chan struct{}) {
	// connect to endpoint
	intConn := connectTCP(c.routeTo.Addr.String())
	extConn := c.conn
	defer intConn.Close()
	defer extConn.Close()
	defer shutdown(ch)

	go forward(extConn, intConn)
	go forward(intConn, extConn)
}

func listenForConnections(addr *net.TCPAddr, r chan connection) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// listen for tcp connections
		conn, _ := listener.Accept()
		route <- connection{
			conn: &conn,
		}
	}
}

func forward(sender, receiver net.Conn) {
	io.Copy(sender, receiver)
}

func shutdown(ch chan struct{}) {
	ch <- struct{}{}
}

func connectTCP(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
