package config

// ServerConfig is the overarching configuration for 'volchestrator server'
type ServerConfig struct {
	Listen ListenConfig `hcl:"listen,block"`
}

// ListenConfig specifies how the server should listen for gRPC requests
type ListenConfig struct {
	Address string `hcl:"address"`
}
