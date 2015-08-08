package trapeze

import (
	"fmt"
	"io"
	"log"
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
	Endpoints    []ServiceEndpoint
	LoadBalancer LoadBalancer
	Connections  map[ServiceEndpoint]map[*net.Addr]Connection
}

func (p *Proxy) Layer4() {
	route := make(chan Connection)
	errCh := make(chan error)

	go p.listenForEndpointRPCs()

	// listen for incoming connections
	// schedule ServiceEndpoints with NextEndpoint
	go listenForConnections(p.Addr, route, errCh)
	for {
		fmt.Println("waiting for conns")
		select {
		case c := <-route:
			fmt.Println("connection received. Routing.")
			p.addConn(c)
			go p.routeConn(c)
		case err := <-errCh:
			fmt.Println("shutting down")
			log.Fatal(err)
		}
	}
}

// listens on p.rpcAddr for join pool or leave pool
// connections
func (p *Proxy) listenForEndpointRPCs() {
	for {

	}
}

func (p *Proxy) addConn(conn Connection) {
	conn = p.LoadBalancer.NextEndpoint(p.Endpoints, conn)
	conn.RouteTo.Connections += 1
}

func (p *Proxy) routeConn(c Connection) {
	// connect to endpoint
	c.CloseCh = make(chan struct{})
	intConn, err := connectTCP(c.RouteTo.Addr.String())
	if err != nil {
		return
	}
	extConn := c.Conn
	defer p.removeConn(c)

	forward(*extConn, intConn)
}

func listenForConnections(addr *net.TCPAddr, r chan Connection, sdCh chan error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		sdCh <- err
	}
	for {
		// listen for tcp connections
		conn, _ := listener.Accept()
		r <- Connection{
			Conn: &conn,
		}
	}
}

// forwards traffic both ways and
func forward(sender, receiver net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	defer sender.Close()
	defer receiver.Close()
	go copyio(receiver, sender, wg)
	go copyio(sender, receiver, wg)
	wg.Wait()
}

func copyio(sender, receiver net.Conn, wg sync.WaitGroup) {
	defer wg.Done()
	io.Copy(sender, receiver)
}

func (p *Proxy) removeConn(conn Connection) {
	// delete: map = p.connections[routeTo] key = conn.conn.RemoteAddr()
	key := *conn.Conn
	keyAddr := key.RemoteAddr()
	delete(p.Connections[*conn.RouteTo], &keyAddr)
	conn.RouteTo.Connections -= 1
	conn.CloseCh <- struct{}{}
}

func connectTCP(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
