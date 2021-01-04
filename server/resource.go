package server

// ResourceManager is responsible for managing the underlying resource represented by a
// Volume, with a given client
type ResourceManager interface {
	Associate(*Volume) error
	Disassociate(*Volume) error
}
