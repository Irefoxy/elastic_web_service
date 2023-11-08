package repo

import "elastic_web_service/internal/model"

type Repo interface {
	GetRecommendPlaces(lon, lat float64) ([]model.Place, error)
	GetPlaces(offset, limit int) ([]model.Place, int, error)
	Init() error
}
