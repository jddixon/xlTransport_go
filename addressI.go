package transport

// An Address provides enough information to identify an endpoint.
// The information needed depends upon the communications protocol
// used.
type AddressI interface {
	Clone() (AddressI, error)
	Equal(any interface{}) bool
	String() string
}
