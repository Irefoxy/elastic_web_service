package repo

import (
	"bytes"
	sModel "elastic_web_service/internal/model"
	"elastic_web_service/internal/repo/converter"
	"elastic_web_service/internal/repo/model"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
)

var _ Repo = (*Elastic)(nil)

const address = "http://localhost:9200"

type Elastic struct {
	client *elasticsearch.Client
	Index  string
}

func NewRepo() *Elastic {
	return &Elastic{}
}

func (es *Elastic) Init() error {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	es.client = esClient
	es.Index = "places"
	return nil
}

func (es *Elastic) GetPlaces(limit, offset int) ([]sModel.Place, int, error) {
	query := map[string]interface{}{
		"size": limit,
		"from": offset,
	}

	places, total, err := es.search(query)
	if err != nil {
		return nil, 0, err
	}
	var servicePlaces = make([]sModel.Place, 0)
	for _, place := range places {
		servicePlaces = append(servicePlaces, converter.ToServiceFromRepo(place))
	}

	return servicePlaces, total, nil
}

func (es *Elastic) GetRecommendPlaces(lon, lat float64) ([]sModel.Place, error) {
	query := map[string]interface{}{
		"size": 3,
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]float64{
						"lat": lon,
						"lon": lat,
					},
					"order":           "asc",
					"unit":            "km",
					"mode":            "min",
					"distance_type":   "arc",
					"ignore_unmapped": true,
				},
			},
		},
	}

	places, _, err := es.search(query)
	if err != nil {
		return nil, err
	}

	var servicePlaces = make([]sModel.Place, 0)
	for _, place := range places {
		servicePlaces = append(servicePlaces, converter.ToServiceFromRepo(place))
	}

	return servicePlaces, nil
}

func (es *Elastic) search(query map[string]interface{}) ([]model.Place, int, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex(es.Index),
		es.client.Search.WithBody(bytes.NewReader(queryBytes)),
		es.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to read response body: %v", err)
		}
		return nil, 0, fmt.Errorf("failed to get places: %s\n%s", res.Status(), responseBody)
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source model.Place `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	places := make([]model.Place, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		places[i] = hit.Source
	}
	return places, result.Hits.Total.Value, nil
}
