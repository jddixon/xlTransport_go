package transport

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var _ = fmt.Print

const (
	QUAD_PAT    = `(?:\d|\d\d|(?:[01]\d\d|2(?:[01234]\d|5[0-5])))`
	V4_ADDR_PAT = `^` + QUAD_PAT + `\.` + QUAD_PAT + `\.` + QUAD_PAT + `\.` + QUAD_PAT + `$`

	bad_port_number = "not a valid IPv4 port number: "
	bad_ipv4_addr   = "not a valid IPv4 address: "
)

var v4AddrRE *regexp.Regexp

func makeRE() (err error) {
	v4AddrRE, err = regexp.Compile(V4_ADDR_PAT)
	return err
}

// An IPv4 address
type V4Address struct {
	host string
	port string // if it's an int, the default is zero
}

// Verify that a string represents a valid port number (in the range
// 0..65535 inclusive).
func checkPortPart(val string) (err error) {
	var port int
	if port, err = strconv.Atoi(val); err == nil {
		if port >= 256*256 {
			err = errors.New(bad_port_number + val)
		} 
	}
	return
}

// Expect an IPV4 address in the form A.B.C.D:P, where P is the
// port number and the :P is optional.  The :P part must be present.:
//
// Accept ":8080" as a valid address, with an implicit "127.0.0.1" host part.  
// Accept "[::]" is a valid host part, interpreted as 0.0.0.0

func NewV4Address(val string) (addr *V4Address, err error) {
	if v4AddrRE == nil {
		if err = makeRE(); err != nil {
			panic(err)
		}
	}
	var addrPart, portPart string

	val = strings.TrimSpace(val)
	if len(val) == 0 {
		err = EmptyAddrString

	} else if val[0] == ':' {
		// accept an address in the form ":nnnn"
		err = checkPortPart(val[1:]) 
		if err == nil {
			addrPart = "127.0.0.1"
			portPart = val[1:]
		}
	} else if strings.HasPrefix(val, "[::]:") {
		err = checkPortPart(val[5:]) 
		if err == nil {
			addrPart = "0.0.0.0"
			portPart = val[5:]
		}
	} else {
		parts := strings.Split(val, `:`)
		partsCount := len(parts)
		if partsCount == 0 || partsCount > 2 {
			err = errors.New(bad_ipv4_addr + val)
		} else if partsCount == 1 {
			// no colon
			if v4AddrRE.MatchString(val) {
				addrPart = val
			} else {
				err = errors.New(bad_ipv4_addr + val)
			}
		} else {
			// we have a colon
			portPart = parts[1]
			err = checkPortPart(portPart) 
			if err == nil {
				addrPart = parts[0]
				if ! v4AddrRE.MatchString(addrPart) {
					err = errors.New(bad_ipv4_addr + val)
				}
			}
		}
	}
	if err == nil {
		addr = &V4Address{addrPart, portPart}
	}
	return
}
func (a *V4Address) Clone() (AddressI, error) {
	return NewV4Address(a.String()) // .(AddressI)
}
func (a *V4Address) Equal(any interface{}) bool {
	if any == nil {
		return false
	}
	if any == a {
		return true
	}
	switch v := any.(type) {
	case *V4Address:
		_ = v
	default:
		return false
	}
	other := any.(*V4Address)
	return a.host == other.host && a.port == other.port
}
func (a *V4Address) String() string {
	if a.port == "" {
		return a.host
	} else {
		return a.host + ":" + a.port
	}
}
