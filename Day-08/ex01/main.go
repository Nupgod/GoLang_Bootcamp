package main

import (
	"fmt"
	"reflect"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

type AndAnotherOneUnknownPlant struct {
	HeightOfRoot int `unit:"metres"`
	LeafNumber    int
	Name      string `language:"russian"`
}



func describePlant(plant interface{}) {
	v := reflect.ValueOf(plant)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		switch field.Kind() {
			case reflect.Int:
				tag := fieldType.Tag
				if tag != "" {
					fmt.Printf("%s(%s):%d\n", fieldType.Name, tag, field.Int())
				} else {
					fmt.Printf("%s:%d\n", fieldType.Name, field.Int())
				}
			case reflect.String:
				tag := fieldType.Tag
				if tag != "" {
					fmt.Printf("%s(%s):%s\n", fieldType.Name, tag, field.String())
				} else {
					fmt.Printf("%s:%s\n", fieldType.Name, field.String())
				}
			default:
				fmt.Printf("Unsupported type: %s\n", field.Kind())
		}
	}
	fmt.Printf("\n")
}

func main() {
	plant1 := UnknownPlant{
		FlowerType: "rose",
		LeafType:   "lanceolate",
		Color:      10,
	}

	plant2 := AnotherUnknownPlant{
		FlowerColor: 10,
		LeafType:    "lanceolate",
		Height:      15,
	}
	plant3 := AndAnotherOneUnknownPlant{
		HeightOfRoot: 6,
		LeafNumber: 26,
		Name: "Пупавка",
	}

	plants := []interface{}{plant1, plant2, plant3}
	for _, plant := range plants{
		describePlant(plant)
	}
}