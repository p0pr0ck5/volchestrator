package config

// ServerConfig is the overarching configuration for 'volchestrator server'
type ServerConfig struct {
	Listen  ListenConfig  `hcl:"listen,block"`
	Backend BackendConfig `hcl:"backend,block"`
}

// ListenConfig specifies how the server should listen for gRPC requests
type ListenConfig struct {
	Address string `hcl:"address"`
}

// BackendConfig specifies how to store data in the backend
type BackendConfig struct {
	Type string `hcl:"type,label"`
}

// ClientConfig is the overarching configuration for 'volchestrator client'
type ClientConfig struct {
	ServerAddress string         `hcl:"server_address,optional"`
	ClientID      string         `hcl:"client_id,optional"`
	LeaseRequests []LeaseRequest `hcl:"lease_request,block"`
}

// LeaseRequest defines a configuration for a client's desire to lease a volume
type LeaseRequest struct {
	Tag              string `hcl:"tag"`
	AvailabilityZone string `hcl:"az"`
}
