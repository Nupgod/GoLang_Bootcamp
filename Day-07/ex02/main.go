package ex02

import (
 "sort"
)

// MinCoins calculates the minimum number of coins needed to represent a value
// using a greedy algorithm. It iterates through the coins in descending order,
// adding the largest denomination coin possible until the remaining value is zero.
// This function does not guarantee the optimal solution for all cases.
// Parameters:
//   - val: The value to represent.
//   - coins: A slice of coin denominations, assumed to be sorted in descending order.
// Returns:
//   - res: A slice containing the denominations used to represent the value.
func MinCoins(val int, coins []int) []int {
 res := make([]int, 0)
 i := len(coins) - 1
 for i >= 0 {
  for val >= coins[i] {
   val -= coins[i]
   res = append(res, coins[i])
  }
  i -= 1
 }
 return res
}

// MinCoins2 calculates the minimum number of coins needed to represent a value
// using a greedy algorithm with sorted denominations. It first sorts the coins
// in ascending order and then applies the same greedy strategy as MinCoins.
// Sorting the denominations allows for better performance in some cases.
// Parameters:
//   - val: The value to represent.
//   - coins: A slice of coin denominations, assumed to be unsorted.
// Returns:
//   - res: A slice containing the denominations used to represent the value.
func MinCoins2(val int, coins []int) []int {
 sort.Ints(coins) // Sort denominations in ascending order
 res := make([]int, 0)
 i := len(coins) - 1
 for i >= 0 {
  for val >= coins[i] {
   val -= coins[i]
   res = append(res, coins[i])
  }
  i -= 1
 }
 return res
}

// To generate HTML documentation:
// 1. Install godoc tool: go get golang.org/x/tools/cmd/godoc
// 2. Run godoc server: godoc -http=:6060
// 3. Access documentation in your browser: http://localhost:6060/pkg/