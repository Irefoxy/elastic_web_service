package converter

import (
	"elastic_web_service/internal/model"
	rModel "elastic_web_service/internal/repo/model"
)

func ToServiceFromRepo(place rModel.Place) model.Place {
	return model.Place{
		Name:    place.Name,
		Address: place.Address,
		Phone:   place.Phone,
		Location: struct {
			Lon float64
			Lat float64
		}{Lon: place.Location.Lon, Lat: place.Location.Lat},
	}
}
