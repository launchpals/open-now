package env

import "os"

// Values denotes various env configuration values
type Values struct {
	GCPKey string
}

// Load loads up all relevant env values
func Load() *Values {
	return &Values{
		GCPKey: os.Getenv("GCP_KEY"),
	}
}
