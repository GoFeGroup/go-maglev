package loadbalancer

// BackendID is the backend's ID.
type BackendID uint32

// ID is the ID of L3n4Addr endpoint (either service or backend).
type ID uint32

// BackendState is the state of a backend for load-balancing service traffic.
type BackendState uint8

// Preferred indicates if this backend is preferred to be load balanced.
type Preferred bool

// Backend represents load balancer backend.
type Backend struct {
	// FEPortName is the frontend port name. This is used to filter backends sending to EDS.
	FEPortName string
	// ID of the backend
	ID BackendID
	// Weight of backend
	Weight uint16
	// Node hosting this backend. This is used to determine backends local to
	// a node.
	NodeName string
	// State of the backend for load-balancing service traffic
	State BackendState
	// Preferred indicates if the healthy backend is preferred
	Preferred Preferred
}
