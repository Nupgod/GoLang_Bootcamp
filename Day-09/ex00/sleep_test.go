package ex00

import (
	"sort"
	"testing"
)

func TestSleepSort(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	expected := make([]int, len(nums))
	copy(expected, nums)
	sort.Ints(expected)

	out := sleepSort(nums)
	actual := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		actual = append(actual, <-out)
	}

	if len(actual) != len(expected) {
		t.Errorf("Expected length %d, but got %d", len(expected), len(actual))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected, actual)
			break
		}
	}
}

func TestEmptySlice(t *testing.T) {
	nums := []int{}
	expected := make([]int, len(nums))
	copy(expected, nums)
	sort.Ints(expected)

	out := sleepSort(nums)
	actual := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		actual = append(actual, <-out)
	}

	if len(actual) != len(expected) {
		t.Errorf("Expected length %d, but got %d", len(expected), len(actual))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected, actual)
			break
		}
	}
}