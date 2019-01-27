package maps

import (
	"go.uber.org/zap"
	gmaps "googlemaps.github.io/maps"
)

// Client interacts with map services
type Client struct {
	l  *zap.SugaredLogger
	gm *gmaps.Client
}

// NewClient instantiates a maps client
func NewClient(l *zap.SugaredLogger, key string) (*Client, error) {
	gm, err := gmaps.NewClient(gmaps.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	return &Client{
		l:  l,
		gm: gm,
	}, nil
}
