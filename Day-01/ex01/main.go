package main

import (
	"flag"
	"fmt"
	"os"
	"src/reader"
	"src/compareDB"
)

func main() {
	oldDatabasePtr := flag.String("old", "", "Original database file to read")
	newDatabasePtr := flag.String("new", "", "Stolen database file to read")
	flag.Parse()

	if *oldDatabasePtr == "" || *newDatabasePtr == "" {
		fmt.Println("Usage: ./compareDB --old <original_database> --new <stolen_database>")
		os.Exit(1)
	}

	// Read the original database
	Reader, err := reader.ChooseDBReader(*oldDatabasePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	originalDB, err := Reader.Read(*oldDatabasePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read the stolen database
	Reader, err = reader.ChooseDBReader(*newDatabasePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stolenDB, err := Reader.Read(*newDatabasePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Compare the databases
	compareDB.CompareDatabases(originalDB, stolenDB)
}
