package ex02

import (
    "sort"
    "sync"
    "testing"
)

func TestMultiplex(t *testing.T) {
    // Test case 1: Test with multiple input channels
    input1 := make(chan interface{})
    input2 := make(chan interface{})
    input3 := make(chan interface{})

    output := multiplex(input1, input2, input3)

    var wg sync.WaitGroup
    wg.Add(3)

    // Write values to input channels
    go func() {
        defer wg.Done()
        for i := 0; i < 3; i++ {
            input1 <- i
        }
        close(input1)
    }()
    go func() {
        defer wg.Done()
        for i := 3; i < 6; i++ {
            input2 <- i
        }
        close(input2)
    }()
    go func() {
        defer wg.Done()
        for i := 6; i < 9; i++ {
            input3 <- i
        }
        close(input3)
    }()

    go func() {
        wg.Wait()
    }()

    // Collect all values in a slice
    var results []interface{}
    for v := range output {
        results = append(results, v)
    }

    // Wait for all values to be received
    sort.Slice(results, func(i, j int) bool {
        return results[i].(int) < results[j].(int)
    })

    // Check if all values are received
    expected := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8}
    if len(results) != len(expected) {
        t.Errorf("Expected %d values, but got %d", len(expected), len(results))
    }

    // Check if all values are correct
    for i, v := range expected {
        if v != results[i] {
            t.Errorf("Expected value %v at index %d, but got %v", v, i, results[i])
        }
    }
}