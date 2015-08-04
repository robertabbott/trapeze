package trapeze

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// implements OSI layer 4 load balancer
// layer four load balancing occurs at the tcp/udp
// level. If you app/site/whatever has multiple
// external IPs a layer 4 load balancer will route
// traffic between the 3 IPs without inspecting packets
// for application specific information

// reverse proxy takes internet traffic and routes in to
// members of the routing pool. The 'routing pool' in this
// case is a slice of ServiceEndpoints

type Proxy struct {
	Addr *net.TCPAddr

	rpcAddr      *net.TCPAddr
	endpoints    []ServiceEndpoint
	loadBalancer LoadBalancer
	connections  map[ServiceEndpoint]map[net.TCPAddr]connection
}

func (p *Proxy) Layer4() {
	route := make(chan connection)

	go p.listenForEndpointRPCs()

	// listen for incoming connections
	// schedule ServiceEndpoints with NextEndpoint
	go listenForConnections(p.Addr, route)
	for {
		select {
		case c := <-route:
			p.addConn(c)
			go routeConn(conn)
		}
	}
}

// listens on p.rpcAddr for join pool or leave pool
// connections
func (p *Proxy) listenForEndpointRPCs() {
	for {

	}
}

func (p *Proxy) addConn(conn connection) {
	conn := p.loadBalancer.NextEndpoint(p.endpoints, c)
	conn.routeTo.connections += 1
}

func (p *Proxy) routeConn(c connection) {
	// connect to endpoint
	c.closeCh = make(chan struct{})
	intConn := connectTCP(c.routeTo.Addr.String())
	extConn := c.conn
	defer p.removeConn(c)

	forward(extConn, intConn)
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

// forwards traffic both ways and
func forward(sender, receiver net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	defer intConn.Close()
	defer extConn.Close()
	go copyio(receiver, sender)
	go copyio(sender, receiver)
	wg.Wait()
}

func copyio(sender, receiver net.Conn, wg sync.WaitGroup) {
	defer wg.Done()
	io.Copy(sender, receiver)
}

func (p *Proxy) removeConn(conn connection) {
	// delete: map = p.connections[routeTo] key = conn.conn.RemoteAddr()
	delete(p.connections[conn.routeTo], conn.conn.RemoteAddr())
	conn.routeTo.connections -= 1
	conn.closeCh <- struct{}{}
}

func connectTCP(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
