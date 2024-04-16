package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	oldSnapshotPtr := flag.String("old", "", "Original snapshot file to read")
	newSnapshotPtr := flag.String("new", "", "New snapshot file to read")
	flag.Parse()

	if *oldSnapshotPtr == "" || *newSnapshotPtr == "" {
		fmt.Println("Usage: ./compareFS --old <original_snapshot> --new <new_snapshot>")
		os.Exit(1)
	}

	// Read the original snapshot and store the lines in a map
	originalLines := make(map[string]bool)
	err := readSnapshotIntoMap(*oldSnapshotPtr, originalLines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read the new snapshot and compare the lines
	err = compareSnapshots(*newSnapshotPtr, originalLines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// readSnapshotIntoMap reads a snapshot file into a map where the keys are the file paths.
func readSnapshotIntoMap(filename string, lines map[string]bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines[line] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// compareSnapshots compares a new snapshot with the original lines and prints the differences.
func compareSnapshots(filename string, originalLines map[string]bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := originalLines[line]; !exists {
			fmt.Printf("ADDED %s\n", line)
		} else {
			delete(originalLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Any remaining lines in originalLines map are files that were removed.
	for removedFile := range originalLines {
		fmt.Printf("REMOVED %s\n", removedFile)
	}

	return nil
}