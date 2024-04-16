package ex01

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestCrawlWeb tests the crawlWeb function with a mock HTTP server.
func TestCrawlWeb(t *testing.T) {
	// Create a mock HTTP server to simulate the web pages.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/page1":
			io.WriteString(w, "This is page 1")
		case "/page2":
			io.WriteString(w, "This is page 2")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create a context for the test.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channels for URLs and done signal.
	urls := make(chan string)
	done := make(chan struct{})

	// Start the crawlWeb function.
	results := crawlWeb(ctx, urls, done)

	// Send URLs to the crawlWeb function.
	go func() {
		urls <- server.URL + "/page1"
		urls <- server.URL + "/page2"
		close(urls)
	}()

	// Collect the results.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range results {
			log.Printf("Received result: %s\n", result)
			// Add assertions here to check the result content.
		}
	}()

	// Wait for the crawling to finish.
	wg.Wait()

	// Cancel the context to stop the crawling.
	cancel()
}
