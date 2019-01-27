package maps

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"googlemaps.github.io/maps"
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
	l.Info("attempting to connect to gmaps")
	if _, _, err := gm.Directions(context.Background(), &maps.DirectionsRequest{
		Origin:      "Vancouver",
		Destination: "Surrey",
	}); err != nil {
		l.Errorw("failed to connect to google maps", "error", err)
		return nil, fmt.Errorf("failed to connect to google maps: %s", err.Error())
	}
	l.Info("successfully made query to gmaps")
	return &Client{
		l:  l,
		gm: gm,
	}, nil
}
