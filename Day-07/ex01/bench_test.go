package ex01
 // go test -bench=. -benchmem -cpuprofile cpu.pprof
 // go tool pprof ex01.test cpu.pprof
import (
	"math/rand"
	"testing"
)

// Generate a slice of random denominations for testing.
func generateDenominations(n int) []int {
	denominations := make([]int, n)
	for i := 0; i < n; i++ {
		denominations[i] = rand.Intn(1000) + 1 // Ensure no zero or negative values.
	}
	return denominations
}

// BenchmarkMinCoins2 benchmarks the minCoins2 function.
func BenchmarkMinCoins2(b *testing.B) {
	// Generate a large slice of denominations for the benchmark.
	denominations := generateDenominations(1000)
	value := rand.Intn(10000) + 1 // Ensure the value is within a reasonable range.
	b.ResetTimer()
	b.ReportAllocs()
	// Run the benchmark.
	for n := 0; n < b.N; n++ {
		MinCoins2(value, denominations)
	}
}
func BenchmarkMinCoins(b *testing.B) {
	// Generate a large slice of denominations for the benchmark.
	denominations := generateDenominations(1000)
	value := rand.Intn(10000) + 1 // Ensure the value is within a reasonable range.
	b.ResetTimer()
	b.ReportAllocs()
	// Run the benchmark.
	for n := 0; n < b.N; n++ {
		MinCoins(value, denominations)
	}
}
