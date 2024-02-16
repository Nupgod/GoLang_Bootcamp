package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Ingredient represents an ingredient in a cake.
type Ingredient struct {
	Name  string  `json:"ingredient_name" xml:"itemname"`
	Count float64 `json:"ingredient_count" xml:"itemcount"`
	Unit  string  `json:"ingredient_unit,omitempty" xml:"itemunit,omitempty"`
}

// Cake represents a cake with its ingredients.
type Cake struct {
	Name         string       `json:"name" xml:"name"`
	Time         string       `json:"time" xml:"stovetime"`
	Ingredients  []Ingredient `json:"ingredients,omitempty" xml:"ingredients>item"`
}

// Recipe represents a collection of cakes.
type Recipe struct {
	Cakes []Cake `json:"cake" xml:"cake"`
}

// DBReader is the interface for reading databases.
type DBReader interface {
	Read(filePath string) (interface{}, error)
}

// JSONDBReader is an implementation of DBReader for JSON files.
type JSONDBReader struct{}

func (r JSONDBReader) Read(filePath string) (interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// XMLDBReader is an implementation of DBReader for XML files.
type XMLDBReader struct{}

func (r XMLDBReader) Read(filePath string) (interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getDBReader returns the appropriate DBReader based on the file extension.
func getDBReader(filePath string) DBReader {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".json":
		return JSONDBReader{}
	case ".xml":
		return XMLDBReader{}
	default:
		panic("Unsupported file format")
	}
}

func main() {
	filePath := flag.String("f", "", "Path to the database file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path using the -f flag.")
		os.Exit(1)
	}

	reader := getDBReader(*filePath)
	data, err := reader.Read(*filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", *filePath, err)
		os.Exit(1)
	}

	// Convert the data to JSON and print it.
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Error converting data to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}