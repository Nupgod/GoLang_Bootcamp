package main

import (
	"flag"
	"fmt"
	"os"
	"src/reader"
)

func main() {
	filenamePtr := flag.String("f", "", "Database file to read")
	flag.Parse()

	if *filenamePtr == "" {
		fmt.Println("Usage: ./readDB -f <filename>")
		os.Exit(1)
	}
	Reader, err := reader.ChooseDBReader(*filenamePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := Reader.Read(*filenamePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Convert the data and print
	Data, err := reader.ConvertData(db, *filenamePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(Data))
}
