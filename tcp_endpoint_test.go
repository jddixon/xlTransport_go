package transport

// xlTransport_go/tcp_endpoint_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestEndPointInterface(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_END_POINT_INTERFACE")
	}
	ep, err := NewTcpEndPoint("127.0.0.1:80")
	c.Assert(err, Equals, nil)

	addr := ep.Address()
	c.Assert(addr.String(), Equals, "127.0.0.1:80")

	x, err := ep.Clone()
	c.Assert(err, Equals, nil)
	c.Assert(ep.String(), Equals, x.String())
	c.Assert(ep.Equal(x), Equals, true)

	c.Assert(ep.Transport(), Equals, "tcp")

	foo := EndPointI(ep) // compiler accepts
	// bar := EndPointI(*ep)		// compiler rejects

	_ = foo
	// _,_ = foo, bar
}

func (s *XLSuite) TestEndPointAnyInterface(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_END_POINT_ANY_INTERFACE")
	}
	ep, err := NewTcpEndPoint("[::]:80")
	c.Assert(err, Equals, nil)

	addr := ep.Address()
	c.Assert(addr.String(), Equals, "0.0.0.0:80")	
	//c.Assert(addr.String(), Equals, "[::]:80")	

	x, err := ep.Clone()
	c.Assert(err, Equals, nil)
	c.Assert(x, NotNil)

	xAddr := x.Address()
	c.Assert(xAddr.String(), Equals, addr.String())
	
	c.Assert(ep.Equal(x), Equals, true)

	c.Assert(ep.Transport(), Equals, "tcp")

	foo := EndPointI(ep) // compiler accepts
	_ = foo
}
