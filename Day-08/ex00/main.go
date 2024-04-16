package main

import (
	"errors"
	"fmt"
	"time"
	"math/rand"
	"runtime"
	"reflect"
	"strings"
)

func getElementBinary(arr []int, idx int) (int, error) {
	// Check for empty slice
	if len(arr) == 0 {
		return 0, errors.New("empty slice")
	}
	// Check for negative index
	if idx < 0 {
		return 0, errors.New("negative index")
	}
	// Check if index is out of bounds
	if idx >= len(arr) {
		return 0, errors.New("index out of bounds")
	}
	// Perform binary search
	low, high := 0, len(arr)-1
	for low <= high {
		mid := (low + high) / 2
		if mid == idx {
			return arr[mid], nil
		} else if mid < idx {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	// This should never be reached
	return 0, errors.New("unexpected error")
}

func getElement(arr []int, idx int) (int, error) {
	// Check for empty slice
	if len(arr) == 0 {
		return 0, errors.New("empty slice")
	}
	// Check for negative index
	if idx < 0 {
		return 0, errors.New("negative index")
	}
	// Check if index is out of bounds
	if idx >= len(arr) {
		return 0, errors.New("index out of bounds")
	}
	// Iterate through the slice
	for i, num := range arr {
		if i == idx {
			return num, nil
		}
	}
	// This should never be reached
	return 0, errors.New("unexpected error")
}

func timerDecorator(fn GetElementFunc) GetElementFunc {
	return func(arr []int, idx int) (int, error) {
		start := time.Now()
		result, err := fn(arr, idx)
		elapsed := time.Since(start)

		// Get the name of the function using reflection
		pc := reflect.ValueOf(fn).Pointer()
		funcName := runtime.FuncForPC(pc).Name()
		funcName = strings.Trim(funcName, "main.") 

		fmt.Printf("Execution time of function %s is %d nanoseconds\n", funcName, elapsed.Nanoseconds())
		return result, err
	}
}
type GetElementFunc func([]int, int) (int, error)

func main() {

	arr := make([]int, 100_000_000)
	// Populate the slice with random values
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000) // Generate a random number between 0 and 999
	}
	idx := rand.Intn(100_000_000)

	dec1 := timerDecorator(getElementBinary)
	dec2 := timerDecorator(getElement)	

	element, err := dec1(arr, idx)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Element at index", idx, ":", element)
	}

	element, err = dec2(arr, idx)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Element at index", idx, ":", element)
	}
	fmt.Printf("True ans: %d", arr[idx])
}
