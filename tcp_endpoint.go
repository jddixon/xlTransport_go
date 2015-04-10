package transport

import (
	"fmt"
	"net"
	"strings"
)

var _ = fmt.Print // DEBUG

var (
	ANY_TCP_END_POINT *TcpEndPoint
)

func init() {
	ANY_TCP_END_POINT, _ = NewTcpEndPoint("127.0.0.1:0")
}

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 *
 */

type TcpEndPoint struct {
	tcpAddr *net.TCPAddr // IP, Port, Zone
}

func NewTcpEndPoint(addr string) (*TcpEndPoint, error) {
	a, err := net.ResolveTCPAddr("tcp", addr)
	if err == nil {
		return &TcpEndPoint{a}, nil
	} else {
		return nil, err
	}
}

func (e *TcpEndPoint) Address() AddressI {
	// return a copy
	a, _ := NewV4Address(e.tcpAddr.String())
	return a
}

func (e *TcpEndPoint) Clone() (ep EndPointI, err error) {
	return NewTcpEndPoint(e.Address().String())
}

func (e *TcpEndPoint) Equal(any interface{}) bool {
	if any == nil {
		return false
	}
	if any == e {
		return true
	}
	switch v := any.(type) {
	case *TcpEndPoint:
		_ = v
	default:
		return false
	}
	other := any.(*TcpEndPoint)
	t, ot := e.tcpAddr, other.tcpAddr

	ts := t.String()
	ots := ot.String()

	if ts != ots {
		if strings.HasPrefix(ts, "[::]") {
			ts = "0.0.0.0" + ts[4:]
		} else if ts[0] == ':' {
			ts = "127.0.0.1" + ts
		}
		if strings.HasPrefix(ots, "[::]") {
			ots = "0.0.0.0" + ots[4:]
		} else if ots[0] == ':' {
			ots = "127.0.0.1" + ots
		}
		if ts != ots {
			return false
		}
	}
	return t.Port == ot.Port && t.Zone == ot.Zone
}

func (e *TcpEndPoint) String() string {
	return "TcpEndPoint: " + e.tcpAddr.String()
}

func (e *TcpEndPoint) Transport() string {
	return "tcp"
}

// net.Addr interface ///////////////////////////////////////////////

// This is just an alias for Transport
func (e *TcpEndPoint) Network() string {
	return e.Transport()
}

// Shortcut for Go
func (e *TcpEndPoint) GetTcpAddr() *net.TCPAddr {
	return e.tcpAddr
}
