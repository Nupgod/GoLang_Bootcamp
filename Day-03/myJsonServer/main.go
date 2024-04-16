package main

import (
	"log"
	"net/http"
	rec "src/myJsonServer/recommend"
	"src/myJsonServer/JWT"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Initialize the Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Resolve the client info for debugging
	info, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer info.Body.Close()

	// Define the route for the places API
	http.HandleFunc("/api/get_token", jwt.GenerateTokenHandler)
	http.HandleFunc("/api/places", rec.GetPlacesHandler(es))
	http.HandleFunc("/api/recommend", jwt.AuthenticateJWT(rec.GetPlacesHandlerRecommed(es)))

	// Start the server
	log.Println("Server is running on http://127.0.0.1:8888")
	err = http.ListenAndServe("127.0.0.1:8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}