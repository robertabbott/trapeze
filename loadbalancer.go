package trapeze

import (
	"fmt"
	"net"
	"sync"
)

// interface for load balancer

type LoadBalancer interface {
	// NextEndpoint chooses the next server to route a request to
	// This method will use some scheduling algorithm to determine
	// which endpoint
	NextEndpoint(endpoints []ServiceEndpoint, client connection) connection
	AddEndpoint(addr string, port int, endpoints []ServiceEndpoint) ([]ServiceEndpoint, error)
	RemoveEndpoint(endpoint ServiceEndpoint)
}

type ServiceEndpoint struct {
	// addr and port probably wont change so
	// no lock needed to guard them
	Addr *net.TCPAddr
	Port int

	connections int          // count active connections
	connLock    sync.RWMutex // guards conn count on this endpoint

	// maybe somewhere down the line associate ServiceEndpoint
	// with resources allocated to that service in order to do
	// weighted round robin
}

type connection struct {
	conn    *net.Conn
	routeTo *ServiceEndpoint
	closeCh chan struct{}
}

func (se *ServiceEndpoint) String() string {
	return fmt.Sprintf("%s:%d", se.Addr, se.Port)
}
