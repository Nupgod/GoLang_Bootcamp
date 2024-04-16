package db

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

type Place struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type Store interface {
	GetPlaces(limit int, offset int) ([]Place, int, error)
}

type ElasticsearchStore struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchStore(client *elasticsearch.Client, index string) *ElasticsearchStore {
	return &ElasticsearchStore{
		client: client,
		index:  index,
	}
}

func (s *ElasticsearchStore) GetPlaces(limit int, offset int) ([]Place, int, error) {

	resp, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithFrom(offset),
		s.client.Search.WithSize(limit),
		s.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	total := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	places := make([]Place, len(hits))
	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		places[i] = Place{
			Name:    source["name"].(string),
			Address: source["address"].(string),
			Phone:   source["phone"].(string),
		}
	}
	return places, total, nil
}
