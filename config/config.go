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
