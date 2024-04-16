package recommend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

// Place represents a place with an ID, name, address, phone, and location.
type Place struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

// Location represents the geographic coordinates of a place.
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// PlacesResponse represents the response structure for the /api/places endpoint.
type PlacesResponse struct {
	Name      string  `json:"name"`
	Total     int     `json:"total"`
	Places    []Place `json:"places"`
	PrevPage  int     `json:"prev_page,omitempty"`
	NextPage  int     `json:"next_page,omitempty"`
	LastPage  int     `json:"last_page"`
}

// ErrorResponse represents the structure for error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// getPlacesFromElasticsearch retrieves places from Elasticsearch.
func getPlacesFromElasticsearch(es *elasticsearch.Client, from int, size int) ([]Place, int, error) {
	// Construct the search request
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	// Perform the search request
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("places"), // Assuming the places data is indexed in an index named "places"
		es.Search.WithBody(bytes.NewReader(body)),
		es.Search.WithFrom(from),
		es.Search.WithSize(size),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	// Check if the request was successful
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, err
		}
		return nil, 0, fmt.Errorf(fmt.Sprintf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"]))
	}

	// Parse the response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}

	// Extract the places from the response
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	total := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	places := make([]Place, len(hits))
	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		places[i] = Place{
			ID:      hit.(map[string]interface{})["_id"].(string),
			Name:    source["name"].(string),
			Address: source["address"].(string),
			Phone:   source["phone"].(string),
			Location: Location{
				Lat: source["location"].(map[string]interface{})["lat"].(float64),
				Lon: source["location"].(map[string]interface{})["lon"].(float64),
			},
		}
	}

	return places, total, nil
}

// getPlacesHandler is a handler function that returns a list of places in JSON format.
func GetPlacesHandler(es *elasticsearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Parse the 'page' query parameter
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			// Return a 400 Bad Request with an error message
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid 'page' value: '" + pageStr + "'"})
			return
		}

		// Define the page size and calculate the offset
		size := 10 // Number of places per page
		from := (page - 1) * size

		// Retrieve places from Elasticsearch
		places, total, err := getPlacesFromElasticsearch(es, from, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate pagination
		prevPage := page - 1
		nextPage := page + 1
		lastPage := (total + size - 1) / size
		if prevPage < 1 {
			prevPage = 0
		}
		if nextPage > lastPage {
			nextPage = 0
		}

		// Construct the response
		response := PlacesResponse{
			Name:      "Places",
			Total:     total,
			Places:    places,
			PrevPage:  prevPage,
			NextPage:  nextPage,
			LastPage:  lastPage,
		}

		// Encode the response to JSON with pretty-printing and write to the response
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ") // Set the prefix and indentation for pretty-printing
		err = encoder.Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// getPlacesHandler is a handler function that returns a list of the three closest places to the user's location.
func GetPlacesHandlerRecommed(es *elasticsearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Parse the 'lat' and 'lon' query parameters
		latStr := r.URL.Query().Get("lat")
		lonStr := r.URL.Query().Get("lon")
		lat, errLat := strconv.ParseFloat(latStr, 64)
		lon, errLon := strconv.ParseFloat(lonStr, 64)
		if errLat != nil || errLon != nil {
			// Return a 400 Bad Request with an error message
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid 'lat' or 'lon' value"})
			return
		}

		// Define the page size
		size := 3 // We only want the three closest places

		// Retrieve places from Elasticsearch
		places, err := getClosestPlacesFromElasticsearch(es, lat, lon, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Construct the response
		response := RecommendationResponse{
			Name:   "Recommendation",
			Places: places,
		}

		// Encode the response to JSON with pretty-printing and write to the response
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ") // Set the prefix and indentation for pretty-printing
		err = encoder.Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// getClosestPlacesFromElasticsearch retrieves the three closest places from Elasticsearch based on the user's location.
func getClosestPlacesFromElasticsearch(es *elasticsearch.Client, lat float64, lon float64, size int) ([]Place, error) {
	// Construct the search request with sorting by geo distance
	query := map[string]interface{}{
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": lat,
						"lon": lon,
					},
					"order":         "asc",
					"unit":          "km",
					"mode":          "min",
					"distance_type": "arc",
				},
			},
		},
		"size": size, // Only return the three closest places
	}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform the search request
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("places"), // Assuming the places data is indexed in an index named "places"
		es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the request was successful
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(fmt.Sprintf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"]))
	}

	// Parse the response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// Extract the places from the response
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	places := make([]Place, len(hits))
	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		places[i] = Place{
			ID:      hit.(map[string]interface{})["_id"].(string),
			Name:    source["name"].(string),
			Address: source["address"].(string),
			Phone:   source["phone"].(string),
			Location: Location{
				Lat: source["location"].(map[string]interface{})["lat"].(float64),
				Lon: source["location"].(map[string]interface{})["lon"].(float64),
			},
		}
	}

	return places, nil
}

// RecommendationResponse represents the response structure for the /api/recommend endpoint.
type RecommendationResponse struct {
	Name   string  `json:"name"`
	Places []Place `json:"places"`
}