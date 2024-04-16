package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Define command-line flags
	meanFlag := flag.Bool("mean", false, "Display the mean")
	medianFlag := flag.Bool("median", false, "Display the median")
	modeFlag := flag.Bool("mode", false, "Display the mode")
	sdFlag := flag.Bool("sd", false, "Display the standard deviation")

	// Parse the command-line flags
	flag.Parse()

	// Check if any flags have been set
	anyFlagSet := *meanFlag || *medianFlag || *modeFlag || *sdFlag

	// If no flags have been set, set all flags to true to display all metrics
	if !anyFlagSet {
		*meanFlag = true
		*medianFlag = true
		*modeFlag = true
		*sdFlag = true
	}

	scanner := bufio.NewScanner(os.Stdin)
	nums := make([]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Println("Invalid input, please enter an integer number separated by newline.")
			continue
		}
		if int(num) >= -100000 && int(num) <= 100000 {
			nums = append(nums, int(num))
		} else {
			fmt.Println("Invalid input, please enter an integer number between -100000 and 100000.")
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	mean := calculateMean(nums)
	mode := calculateMode(nums)
	median := calculateMedian(nums)
	sd := calculateSD(nums, mean)

	// Display the chosen metrics
	if *meanFlag {
		fmt.Printf("Mean: %.2f\n", mean)
	}
	if *medianFlag {
		fmt.Printf("Median: %.2f\n", median)
	}
	if *modeFlag {
		fmt.Printf("Mode: %v\n", mode)
	}
	if *sdFlag {
		fmt.Printf("SD: %.2f\n", sd)
	}
}


func calculateMean(nums []int) float64 {
	sum := 0.0
	if len(nums) == 0 {return 0}
	for _, num := range nums {
		sum += float64(num)
	}
	return sum / float64(len(nums))
}

func calculateMode(nums []int) int {
	if len(nums) == 0 {return 0}
	counts := make(map[int]int)
	for _, num := range nums {
		counts[num]++
	}
	maxCount := 0
	mode := 0
	for num, count := range counts {
		if count > maxCount {
			maxCount = count
			mode = num
		}
	}
	return mode
}

func calculateMedian(nums []int) float64 {
	length := len(nums)
	if length == 0 {return 0}
	sort.Ints(nums)
	if length%2 == 0 {
		return float64(nums[length/2-1]+nums[length/2]) / 2.0
	} else {
		return float64(nums[length/2])
	}
}

func calculateSD(nums []int, mean float64) float64 {
	if len(nums) == 0 {return 0}
	sum := 0.0
	for _, num := range nums {
		sum += math.Pow(float64(num)-mean, 2)
	}
	return math.Sqrt(sum / float64(len(nums)))
}