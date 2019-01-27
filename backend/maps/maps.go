package maps

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	gmaps "googlemaps.github.io/maps"

	open_now "github.com/launchpals/open-now/proto/go"
)

// Client interacts with map services
type Client struct {
	l  *zap.SugaredLogger
	gm *gmaps.Client

	cache *cache
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
		l:     l,
		gm:    gm,
		cache: newCache(5*time.Minute, 5*time.Minute),
	}, nil
}

// PointsOfInterest returns a series of points
func (c *Client) PointsOfInterest(
	ctx context.Context,
	coords *open_now.Coordinates,
	situation open_now.Context_Situation,
) ([]open_now.Interest, error) {
	var radius uint
	switch situation {
	case open_now.Context_FOOT:
		radius = 1000
	case open_now.Context_VEHICLE:
		radius = 10000
	default:
		radius = 5000
	}

	resp, err := c.gm.TextSearch(ctx, &gmaps.TextSearchRequest{
		Location: &gmaps.LatLng{
			Lat: coords.GetLatitude(),
			Lng: coords.GetLongitude(),
		},
		Radius:  radius,
		OpenNow: true, // gg
	})
	if err != nil {
		c.l.Error("failed to make query", "error", err)
		return nil, err
	}

	pois := []open_now.Interest{}

	for _, l := range resp.Results {
		openTime, err := strconv.ParseInt(l.OpeningHours.Periods[0].Open.Time, 10, 64)

		if err != nil {
			// TODO
		}

		closeTime, err := strconv.ParseInt(l.OpeningHours.Periods[0].Close.Time, 10, 64)

		poi := open_now.Interest{
			InterestId:  l.ID,
			Name:        l.Name,
			Description: "",
			OpeningTime: openTime,
			ClosingTime: closeTime,
			Coordinates: &open_now.Coordinates{
				Latitude:  l.Geometry.Location.Lat,
				Longitude: l.Geometry.Location.Lng,
			},
		}

		pois = append(pois, poi)

		c.l.Debugw("location", "l", l)
	}
	return pois, nil
}

// Close stops background jobs
func (c *Client) Close() { c.cache.stop <- true }
