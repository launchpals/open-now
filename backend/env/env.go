package env

import "os"

// Values denotes various env configuration values
type Values struct {
	Host   string
	GCPKey string
}

// Load loads up all relevant env values
func Load() *Values {
	host := os.Getenv("OPEN_NOW_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	return &Values{
		Host:   os.Getenv("OPEN_NOW_HOST"),
		GCPKey: os.Getenv("GCP_KEY"),
	}
}
