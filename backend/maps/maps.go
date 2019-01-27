package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	gmaps "googlemaps.github.io/maps"

	open_now "github.com/launchpals/open-now/proto/go"
)

// Client interacts with map services
type Client struct {
	l    *zap.SugaredLogger
	gm   *gmaps.Client
	wkey string
}

// NewClient instantiates a maps client
func NewClient(l *zap.SugaredLogger, gkey string, wkey string) (*Client, error) {
	gm, err := gmaps.NewClient(gmaps.WithAPIKey(gkey))
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
		l:    l,
		gm:   gm,
		wkey: wkey,
	}, nil
}

// PointsOfInterest returns a series of points
func (c *Client) PointsOfInterest(
	ctx context.Context,
	coords *open_now.Coordinates,
	situation open_now.Context_Situation,
) ([]*open_now.Interest, error) {
	var radius float64
	switch situation {
	case open_now.Context_FOOT:
		radius = 500
	case open_now.Context_VEHICLE:
		radius = 2000
	default:
		radius = 500
	}

	// TODO: weather call
	urlString := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s",
		coords.GetLatitude(),
		coords.GetLongitude(),
		c.wkey,
	)
	c.l.Debugw("Making weather API call", "target", urlString)

	wResp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}
	defer wResp.Body.Close()

	var raw interface{}
	if err = json.NewDecoder(wResp.Body).Decode(&raw); err != nil {
		c.l.Errorw("failed to read weather response", "error", err)
		return nil, err
	}

	rawJSON := raw.(map[string]interface{})

	if rawJSON["rain"] != nil {
		radius *= 0.75
		c.l.Infof("IT IS RAINING")
	} else {
		c.l.Infof("IT IS NOT RAINING")
	}

	main, _ := rawJSON["main"].(map[string]interface{})
	tempRaw, _ := main["temp"]
	temp, _ := tempRaw.(float64)

	// >35 deg c
	if temp > 308 {
		radius *= 0.75
		c.l.Infof("IT IS TOO HOT: %f", temp)
	} else {
		c.l.Infof("IT IS NOT TOO HOT: %f", temp)
	}

	resp, err := c.gm.NearbySearch(ctx, &gmaps.NearbySearchRequest{
		Location: &gmaps.LatLng{
			Lat: coords.GetLatitude(),
			Lng: coords.GetLongitude(),
		},
		Radius:  uint(radius),
		OpenNow: true, // gg
	})
	if err != nil {
		c.l.Errorw("failed to make query", "error", err)
		return nil, err
	}

	var pois = []*open_now.Interest{}
	for _, l := range resp.Results {
		c.l.Debugw("received response",
			"hours", l.OpeningHours)

		var ot open_now.Interest_Type
		for _, t := range l.Types {
			switch gmaps.PlaceType(t) {
			case gmaps.PlaceTypeEmbassy, gmaps.PlaceTypeCityHall, gmaps.PlaceTypePolice, gmaps.PlaceTypeLocalGovernmentOffice, gmaps.PlaceTypeHospital:
				ot = open_now.Interest_AUTHORITY
				break
			case gmaps.PlaceTypeGroceryOrSupermarket, gmaps.PlaceTypeCafe, gmaps.PlaceTypeBakery, gmaps.PlaceTypeBar, gmaps.PlaceTypeNightClub, gmaps.PlaceTypeFood:
				ot = open_now.Interest_FOOD
				break
			case gmaps.PlaceTypeShoppingMall, gmaps.PlaceTypeStore, gmaps.PlaceTypeEstablishment:
				ot = open_now.Interest_STORE
				break
			case gmaps.PlaceTypeLodging, gmaps.PlaceTypeCampground:
				ot = open_now.Interest_LODGING
				break
			case gmaps.PlaceTypeAmusementPark, gmaps.PlaceTypeAquarium, gmaps.PlaceTypeArtGallery, gmaps.PlaceTypeLibrary:
				ot = open_now.Interest_ATTRACTION
				break
			}
		}
		pois = append(pois, &open_now.Interest{
			InterestId: l.ID,
			Name:       l.Name,
			Type:       ot,
			Photos:     newPhotos(l.Photos),

			LocationDescription: newLocationDescription(&l),
			InterestDescription: strings.Join(l.Types, ", "),
			Coordinates: &open_now.Coordinates{
				Latitude:  l.Geometry.Location.Lat,
				Longitude: l.Geometry.Location.Lng,
			},
		})
	}
	return pois, nil
}

func newPhotos(gphotos []gmaps.Photo) []*open_now.Interest_Photo {
	var photos = make([]*open_now.Interest_Photo, len(gphotos))
	for i := 0; i < len(gphotos); i++ {
		photos[i] = &open_now.Interest_Photo{
			PhotoRef:     gphotos[i].PhotoReference,
			Attributions: gphotos[i].HTMLAttributions,
		}
	}
	return photos
}

func newLocationDescription(l *gmaps.PlacesSearchResult) string {
	if l.FormattedAddress == "" {
		return l.Vicinity
	}
	if l.Vicinity == "" {
		return l.FormattedAddress
	}
	return fmt.Sprintf("%s (near %s)", l.FormattedAddress, l.Vicinity)
}
