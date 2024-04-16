package main

import (
	"fmt"
	"sort"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func getNCoolestPresents(presents []Present, n int) (PresentHeap, error) {
	if n > len(presents) || n < 0 {
		return nil, fmt.Errorf("n is larger than the size of the slice or is negative")
	}
	if len(presents) == 0 {
		return nil, fmt.Errorf("Length of presents must be greater then 0")
	}
	// Sort the slice by Value in descending order, then by Size in ascending order
	sort.Slice(presents, func(i, j int) bool {
		if presents[i].Value == presents[j].Value {
			return presents[i].Size < presents[j].Size
		}
		return presents[i].Value > presents[j].Value
	})
	coolestPresents := make(PresentHeap, 0)

	// Push presents onto the heap, ensuring it contains only the n coolest presents.
	for _, present := range presents {
		if len(coolestPresents) < n {
			coolestPresents = append(coolestPresents, present)
		} else {
			break
		}
	}

	return coolestPresents, nil
}

func main() {
	presents := []Present{
		{5, 1},
		{4, 5},
		{3, 1},
		{5, 2},
	}
	n := 2
	coolestPresents, err := getNCoolestPresents(presents, n)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(coolestPresents)
	}

}
