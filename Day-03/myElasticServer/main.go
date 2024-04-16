package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"src/myElasticServer/db"

	"github.com/elastic/go-elasticsearch/v8"
)

type server struct {
	store db.Store
}

func (s *server) handlePlaces(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 || page > 1365 {
		http.Error(w, "Invalid 'page' value", http.StatusBadRequest)
		return
	}

	limit := 10 // Number of places per page
	offset := (page - 1) * limit
	places, total, err := s.store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastPage := total / limit
	if total%limit != 0 {
		lastPage++
	}

	data := struct {
		Places   []db.Place
		Total    int
		Page     int
		LastPage int
	}{
		Places:   places,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
	}

	// Define custom template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	}

	// Parse the template file with the custom functions
	tmpl, err := template.New("places.html").Funcs(funcMap).ParseFiles(filepath.Join("templates", "places.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	store := db.NewElasticsearchStore(es, "places") // Replace with your Elasticsearch index name
	s := &server{store: store}

	http.HandleFunc("/", s.handlePlaces)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
