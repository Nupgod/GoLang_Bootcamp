package mincoin

import (
	"reflect"
	"testing"
)

func TestMinCoins(t *testing.T) {
	tests := []struct {
		name     string
		val      int
		coins    []int
		expected []int
	}{
		{
			name:     "Test with sorted unique coins",
			val:      13,
			coins:    []int{1, 5, 10},
			expected: []int{10, 1, 1, 1},
		},
		{
			name:     "Test with unsorted coins",
			val:      13,
			coins:    []int{10, 1, 5},
			expected: []int{10, 1, 1, 1},
		},
		{
			name:     "Test with duplicate coins",
			val:      13,
			coins:    []int{1, 5, 5, 10},
			expected: []int{10, 1, 1, 1},
		},
		{
			name:     "Test with empty coins",
			val:      13,
			coins:    []int{},
			expected: []int{},
		},
		{
			name:     "Test with zero value",
			val:      0,
			coins:    []int{1, 5, 10},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minCoins(tt.val, tt.coins)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}