package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

var (
	flagLines    = flag.Bool("l", false, "Count lines")
	flagChars    = flag.Bool("m", false, "Count characters")
	flagWords    = flag.Bool("w", false, "Count words")
)

func main() {
	flag.Parse()
	filenames := flag.Args()

	if len(filenames) == 0 {
		fmt.Println("Please provide at least one filename.")
		os.Exit(1)
	}

	// Ensure that only one flag is specified
	flagsCount := 0
	if *flagLines {
		flagsCount++
	}
	if *flagChars {
		flagsCount++
	}
	if *flagWords {
		flagsCount++
	}
	if flagsCount > 1 {
		fmt.Println("Please specify only one of the following flags: -l, -m, -w")
		os.Exit(1)
	}

	// If no flag is specified, default to -w (word count)
	if flagsCount == 0 {
		*flagWords = true
	}

	var wg sync.WaitGroup

	for _, filename := range filenames {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			countStats(filename)
		}(filename)
	}

	wg.Wait()
}

func countStats(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %q: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	chars := 0
	words := 0

	for scanner.Scan() {
		line := scanner.Text()
		lines++
		chars += utf8.RuneCountInString(line)
		words += countWords(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %q: %v\n", filename, err)
		return
	}

	if *flagLines {
		fmt.Printf("%d\t%s\n", lines, filename)
	} else if *flagChars {
		fmt.Printf("%d\t%s\n", chars, filename)
	} else if *flagWords {
		fmt.Printf("%d\t%s\n", words, filename)
	}
}

func countWords(line string) int {
	// Split the line into words using strings.Fields, which splits on whitespace.
	words := strings.Fields(line)
	return len(words)
}