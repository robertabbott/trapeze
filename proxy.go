package trapeze

import (
	"fmt"
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
}

// keep slice of connections

func (p *Proxy) Layer3() {
	route := make(chan connection)

	// listen for incoming connections
	// schedule ServiceEndpoints based on some algorithm
	// pass scheduling algo to Layer3?

	go listenForConnections(p.Addr, route)
	for {
		select {
		case c := <-route:
			conn := p.loadBalancer.NextEndpoint(p.Endpoints, c)
			go routeConn(conn)
		}
	}
}

func routeConn(c connection) error {
	// do routing

	return nil
}

func listenForConnections(addr *net.TCPAddr, r chan connection) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// listen for tcp connections
	}
}
