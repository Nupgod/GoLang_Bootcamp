package reader

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
)

// IngredientCount is a custom type that can unmarshal from a JSON string or number.
type IngredientCount float64

func (c *IngredientCount) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*c = IngredientCount(value)
	case string:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		*c = IngredientCount(parsed)
	default:
		return fmt.Errorf("unsupported type for IngredientCount: %T", value)
	}
	return nil
}

// Item represents a generic item with a name, count, and unit.
type Item struct {
	Name  string          `xml:"itemname" json:"ingredient_name"`
	Count IngredientCount `xml:"itemcount" json:"ingredient_count"`
	Unit  string          `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

// Cake represents a cake with a name, stovetime, and ingredients.
type Cake struct {
	Name        string `xml:"name" json:"name"`
	StoveTime   string `xml:"stovetime" json:"time"`
	Ingredients []Item `xml:"ingredients>item" json:"ingredients"`
}

// Recipe represents a list of cakes.
type Recipe struct {
	Cakes []Cake `xml:"cake" json:"cake"`
}

// DBReader is the interface for reading databases.
type DBReader interface {
	Read(filename string) (*Recipe, error)
}

// XMLReader is an implementation of DBReader for XML files.
type XMLReader struct{}

func (x *XMLReader) Read(filename string) (*Recipe, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var db Recipe
	err = xml.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// JSONReader is an implementation of DBReader for JSON files.
type JSONReader struct{}

func (j *JSONReader) Read(filename string) (*Recipe, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var db Recipe
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// convertData converts data between JSON and XML formats.
func ConvertData(data *Recipe, filename string) ([]byte, error) {
	ext := filepath.Ext(filename)
	if ext == ".json" {
		return xml.MarshalIndent(data, "", "    ")
	} else {
		return json.MarshalIndent(data, "", "    ")
	}
}


func ChooseDBReader(filename string) (DBReader, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".xml":
		return &XMLReader{}, nil
	case ".json":
		return &JSONReader{}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}