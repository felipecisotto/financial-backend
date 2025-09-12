package mcp

type Config struct {
	ServerName    string
	ServerVersion string
	Transport     TransportType
}

type TransportType string

const (
	TransportTypeStdio TransportType = "stdio"
	TransportTypeHTTP  TransportType = "http"
)

func DefaultConfig() *Config {
	return &Config{
		ServerName:    "financial-backend",
		ServerVersion: "v1.0.0",
		Transport:     TransportTypeStdio,
	}
}
