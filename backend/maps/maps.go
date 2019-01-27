package maps

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	gmaps "googlemaps.github.io/maps"

	open_now "github.com/launchpals/open-now/proto/go"
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

// PointsOfInterest returns a series of points
func (c *Client) PointsOfInterest(
	ctx context.Context,
	coords *open_now.Coordinates,
	situation open_now.Context_Situation,
) ([]*open_now.Interest, error) {
	var radius uint
	switch situation {
	case open_now.Context_FOOT:
		radius = 1000
	case open_now.Context_VEHICLE:
		radius = 10000
	default:
		radius = 5000
	}
	resp, err := c.gm.NearbySearch(ctx, &gmaps.NearbySearchRequest{
		Location: &gmaps.LatLng{
			Lat: coords.GetLatitude(),
			Lng: coords.GetLongitude(),
		},
		Radius:  radius,
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

// Close stops background jobs
func (c *Client) Close() { c.cache.stop <- true }
