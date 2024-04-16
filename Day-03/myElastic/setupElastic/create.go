package setupelastic

import(
	"time"
	"log"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
)

func CreateIndex() (*elasticsearch.Client, error) {
	retryBackoff := backoff.NewExponentialBackOff()
	// Initialize a new Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		// Retry on 429 TooManyRequests statuses
		//
		RetryOnStatus: []int{502, 503, 504, 429},

		// Configure the backoff function
		//
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},

		// Retry up to 5 attempts
		//
		MaxRetries: 5,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		return nil, err
	}
	// Define the index name
	indexName := "places"

	// Check if the index exists
	exists, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		log.Fatalf("Error checking index existence: %s", err)
	}
	exists.Body.Close()
	settings := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {		  
				"name": { "type": "text" },
				"address": { "type": "text" },
				"phone": { "type": "text" },
				"location": { "type": "geo_point" }
			}
		}
	}`
	// If the index does not exist, create it with settings and mappings
	if exists.StatusCode == 404 {
		res, err := es.Indices.Create(indexName,
			es.Indices.Create.WithBody(strings.NewReader(settings)),
		)
		if err != nil {
			log.Fatalf("Error creating the index: %s", err)
			return nil, err
		}
		if res.IsError() {
			log.Fatalf("Error response: %s", res.String())
		} else {
			fmt.Println("Index created successfully")
		}
		defer res.Body.Close()
	} else if exists.StatusCode == 200 {
		fmt.Println("Index already exists")
		if res, err := es.Indices.Delete([]string{indexName}, es.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
			log.Fatalf("Cannot delete index: %s", err)
			return nil, err
		} else {
			res, err := es.Indices.Create(indexName,
				es.Indices.Create.WithBody(strings.NewReader(settings)),
			)
			if err != nil {
				log.Fatalf("Error creating the index: %s", err)
				return nil, err
			}
			if res.IsError() {
				log.Fatalf("Error response: %s", res.String())
			} else {
				fmt.Println("Index recreated successfully")
			}
		}
	} else {
		log.Fatalf("Unexpected response code: %d", exists.StatusCode)
		return nil, err
	}
	return es, nil
}
