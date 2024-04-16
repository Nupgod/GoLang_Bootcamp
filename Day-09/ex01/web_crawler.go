package ex01

import (
    "context"
    "io"
    "log"
    "net/http"
    "sync"
)

// crawlWeb accepts an input channel for sending URLs and returns another
// channel for crawling results.
func crawlWeb(ctx context.Context, urls <-chan string, done <-chan struct{}) <-chan string {
    resultChan := make(chan string)

    // Worker function to process URLs concurrently
    worker := func(url string, wg *sync.WaitGroup) {
        defer wg.Done()

        // Create a new HTTP request
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for URL %s: %v\n", url, err)
            return
        }

        // Send the request and get the response
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            log.Printf("Error fetching URL %s: %v\n", url, err)
            return
        }
        defer resp.Body.Close()

        // Read the response body
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading body for URL %s: %v\n", url, err)
            return
        }

        // Send the body content to the result channel
        select {
        case <-ctx.Done():
            log.Printf("Crawling stopped for URL %s\n", url)
        case <-done:
            return
        case resultChan <- string(body):
        }
    }

    // Control the number of workers using a semaphore
    sem := make(chan struct{}, 8)
    var wg sync.WaitGroup

    // Start workers
    go func() {
        defer close(resultChan)
        for url := range urls {
            select {
            case <-ctx.Done():
                return
            case <-done:
                return
            default:
                // Acquire semaphore
                sem <- struct{}{}
                wg.Add(1)
                go func(url string) {
                    defer func() { <-sem }()
                    worker(url, &wg)
                }(url)
            }
        }
        wg.Wait()
    }()

    return resultChan
}

