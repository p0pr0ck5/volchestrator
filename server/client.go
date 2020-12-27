package server

import "time"

// ClientStatus describes client's current status
type ClientStatus int

const (
	// UnknownClientStatus indicates the client status is unknown
	UnknownClientStatus ClientStatus = iota

	// AliveClientStatus indicates the client is alive
	AliveClientStatus

	// DeadClientStatus indicate the client is dead/unresponsive
	DeadClientStatus

	// LeftClientStatus indicates the client intentionally left
	LeftClientStatus
)

// ClientInfo details information about a given client
type ClientInfo struct {
	ID        string
	Status    ClientStatus
	FirstSeen time.Time
	LastSeen  time.Time
}

// ClientFilterFunc is a function to filter a list of clients based on a given condition
type ClientFilterFunc func(ClientInfo) bool

// ClientFilterAll returns all clients
func ClientFilterAll(ci ClientInfo) bool {
	return true
}

// ClientFilterByStatus returns clients that match a given status
func ClientFilterByStatus(status ClientStatus) ClientFilterFunc {
	return func(ci ClientInfo) bool {
		return ci.Status == status
	}
}
