package trapeze

import (
	"fmt"
	"net"
)

// interface for load balancer

type LoadBalancer interface {
	// NextEndpoint chooses the next server to route a request to
	// This method will use some scheduling algorithm to determine
	// which endpoint
	NextEndpoint(endpoints []ServiceEndpoint, client Connection) Connection
	AddEndpoint(addr string, port int, endpoints []ServiceEndpoint) ([]ServiceEndpoint, error)
	RemoveEndpoint(endpoint ServiceEndpoint)
}

type ServiceEndpoint struct {
	// addr and port probably wont change so
	// no lock needed to guard them
	Addr *net.TCPAddr
	Port int

	Connections int // count active Connections

	// maybe somewhere down the line associate ServiceEndpoint
	// with resources allocated to that service in order to do
	// weighted round robin
}

type Connection struct {
	Conn    *net.Conn
	RouteTo *ServiceEndpoint
	CloseCh chan struct{}
}

func (se *ServiceEndpoint) String() string {
	return fmt.Sprintf("%s:%d", se.Addr, se.Port)
}
