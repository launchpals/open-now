package transit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	open_now "github.com/launchpals/open-now/proto/go"
	"go.uber.org/zap"
)

type StopsResp struct {
	Stops []struct {
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		OnestopID        string `json:"onestop_id"`
		Name             string `json:"name"`
		GeometryCentroid struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry_centroid"`
		RoutesServingStop []struct {
			OperatorName      string `json:"operator_name"`
			OperatorOnestopID string `json:"operator_onestop_id"`
			RouteName         string `json:"route_name"`
			RouteOnestopID    string `json:"route_onestop_id"`
		} `json:"routes_serving_stop"`
	} `json:"stops"`
	Meta struct {
		SortKey   string `json:"sort_key"`
		SortOrder string `json:"sort_order"`
		PerPage   int    `json:"per_page"`
		Offset    int    `json:"offset"`
	} `json:"meta"`
}

// Client interacts with map services
type Client struct {
	l *zap.SugaredLogger
}

// NewClient instantiates a maps client
func NewClient(l *zap.SugaredLogger) (*Client, error) {
	return &Client{
		l: l,
	}, nil
}

// TransitStops returns a list of transit stops
func (c *Client) TransitStops(ctx context.Context, coords *open_now.Coordinates) ([]*open_now.TransitStop, error) {
	urlString := fmt.Sprintf(
		"https://transit.land/api/v1/stops?lat=%f&lon=%f&r=%d",
		coords.GetLatitude(),
		coords.GetLongitude(),
		1000,
	)

	c.l.Infof("Making transit API call: %s", urlString)

	resp, err := http.Get(urlString)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	target := StopsResp{}
	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		c.l.Infow("sHJXT", "error", err)
	}

	var stops = []*open_now.TransitStop{}

	for _, stop := range target.Stops {
		var routes = []string{}

		for _, route := range stop.RoutesServingStop {
			routes = append(routes, route.RouteName)
		}

		stops = append(stops, &open_now.TransitStop{
			Coordinates: &open_now.Coordinates{
				Latitude:  stop.GeometryCentroid.Coordinates[0],
				Longitude: stop.GeometryCentroid.Coordinates[1],
			},
			Routes: routes,
		})
	}

	return stops, nil
}
