package server

// Backend defines functions implemented by the data store
type Backend interface {
	AddClient(string) error
	UpdateClient(string, ClientStatus) error
	RemoveClient(string) error
	Clients(ClientFilterFunc) ([]ClientInfo, error)

	GetVolume(id string) (*Volume, error)
	ListVolumes() ([]*Volume, error)
	AddVolume(*Volume) error
	UpdateVolume(*Volume) error
	DeleteVolume(id string) error
}
