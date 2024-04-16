package main

import (
	"log"
	el "src/myElastic/setupElastic"
)

func main() {
	es, err := el.CreateIndex()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}
	el.UploadData(es)
}
