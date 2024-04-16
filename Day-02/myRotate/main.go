package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	archiveFlag = flag.String("a", "", "Archive all logs to the specified directory")
)

func main() {
	flag.Parse()
	logFiles := flag.Args()
	var wg sync.WaitGroup

	for _, logFile := range logFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			if *archiveFlag != "" {
				// Check if the archive directory exists and is accessible
				archivePathInfo, err := os.Stat(*archiveFlag)
				fmt.Println(archivePathInfo)
				if err != nil || !archivePathInfo.IsDir() {
					fmt.Printf("Error: The specified archive directory %s does not exist or is not accessible.\n", *archiveFlag)
					os.Exit(1)
				}
				rotateAndArchiveLog(*archiveFlag, file)
			} else {
				rotateLog(file)
			}
		}(logFile)
	}

	wg.Wait()
}

func rotateAndArchiveLog(archivePath, logFile string) {
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		fmt.Printf("Error getting file info for %s: %v\n", logFile, err)
		return
	}

	timestamp := fileInfo.ModTime().Unix()
	newLogFile := filepath.Join(archivePath, fmt.Sprintf("%s_%d.tag.gz", filepath.Base(logFile), timestamp))

	err = compressFile(logFile, newLogFile)
	if err != nil {
		fmt.Printf("Error compressing %s: %v\n", logFile, err)
		return
	}

	fmt.Printf("Compressed %s to %s\n", logFile, newLogFile)
}

func rotateLog(logFile string) {
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		fmt.Printf("Error getting file info for %s: %v\n", logFile, err)
		return
	}

	timestamp := fileInfo.ModTime().Unix()
	newLogFile := fmt.Sprintf("%s_%d.tag.gz", logFile, timestamp)

	err = compressFile(logFile, newLogFile)
	if err != nil {
		fmt.Printf("Error compressing %s: %v\n", logFile, err)
		return
	}

	fmt.Printf("Compressed %s to %s\n", logFile, newLogFile)
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		return err
	}

	return nil
}