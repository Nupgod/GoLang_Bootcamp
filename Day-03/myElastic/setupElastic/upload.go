package setupelastic

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func UploadData(es *elasticsearch.Client) {
	// Define a struct to map each row of the CSV file to
	type Document struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Phone    string `json:"phone"`
		Location struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"location"`
	}
	// Open the CSV file
	csvFilePath := filepath.Join("..", "..", "materials", "data.csv")
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %s", err)
	}
	defer csvFile.Close()

	// Create a new CSV reader
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t' // Assuming the CSV is tab-delimited

	// Read the header row
	_, err = reader.Read() // Skip the header
	if err != nil {
		log.Fatalf("Error reading header row: %s", err)
	}

	// Prepare the Bulk API payload
	var bulkPayload strings.Builder

	// Read the rest of the CSV file
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error reading CSV line: %s", err)
		}

		// Parse the latitude and longitude values
		lat, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			log.Fatalf("Error parsing latitude: %s", err)
		}
		lon, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			log.Fatalf("Error parsing longitude: %s", err)
		}
		id, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			log.Fatalf("Error parsing id: %s", err)
		}
		// Create a Document struct and populate it with the CSV data
		doc := Document{
			Name:    line[1],
			Address: line[2],
			Phone:   line[3],
			Location: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{
				Lat: lat,
				Lon: lon,
			},
		}
		docJSON, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Error marshaling document to JSON: %s", err)
		}

		// Add the action and document to the Bulk API payload
		meta := fmt.Sprintf(`{"index":{"_index":"%s", "_id":"%d"}}%s`, "places", id, "\n")
		bulkPayload.WriteString(meta)
		bulkPayload.Write(docJSON)
		bulkPayload.Write([]byte("\n"))
	}

	// Perform the Bulk API request
	res, err := es.Bulk(strings.NewReader(bulkPayload.String()), es.Bulk.WithIndex("places"))
	if err != nil {
		log.Fatalf("Error performing Bulk API request: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response from Elasticsearch: %s", res.String())
	} else {
		fmt.Println("Data loaded successfully into Elasticsearch using Bulk API.")
	}

}
