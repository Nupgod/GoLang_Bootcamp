package main

import (
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

// grabPresents uses dynamic programming to solve the Knapsack problem.
// It returns the presents that maximize the total value while fitting within the given capacity.
func grabPresents(presents []Present, capacity int) []Present {
	n := len(presents)

	// dp[i][w] will store the maximum value that can be achieved
	// with a capacity of 'w' and using only the first 'i' presents.
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity + 1)
	}

	// Fill dp[][] in a bottom-up manner.
	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			if presents[i-1].Size <= w {
				// If the present can fit in the current capacity, consider both
				// including and excluding the present to maximize the value.
				dp[i][w] = max(presents[i-1].Value+dp[i-1][w-presents[i-1].Size], dp[i-1][w])
			} else {
				// If the present is too big for the current capacity, exclude it.
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	// Reconstruct the presents that were used to achieve the optimal value.
	var selectedPresents []Present
	w := capacity
	for i := n; i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			// This present was used.
			selectedPresents = append(selectedPresents, presents[i-1])
			w -= presents[i-1].Size
		}
	}

	// Reverse the selected presents to get them in the order they were given.
	for i, j := 0, len(selectedPresents)-1; i < j; i, j = i+1, j-1 {
		selectedPresents[i], selectedPresents[j] = selectedPresents[j], selectedPresents[i]
	}

	return selectedPresents
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sum(presents []Present) (int, int) {
	var sumVal int
	var capVal int
	for _, present := range presents {
		sumVal += present.Value
		capVal += present.Size
	}
	return sumVal, capVal
}

func main() {
	presents := []Present{
		{60, 10},
		{100, 20},
		{120, 30},
		{75,25},
		{38,7},
		{83,26},
	}
	capacity := 83

	selectedPresents := grabPresents(presents, capacity)
	sum, cap := sum(selectedPresents)
	fmt.Printf("Total value: %v\nTotal size: %v\n", sum, cap)
	fmt.Println("Selected Presents:")
	for _, present := range selectedPresents {
		fmt.Printf("Value: %d, Size: %d\n", present.Value, present.Size)
	}
}