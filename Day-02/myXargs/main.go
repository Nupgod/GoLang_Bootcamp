package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myXargs command [args...]")
		os.Exit(1)
	}

	// The command and its arguments are the first arguments passed to the program.
	command := os.Args[1]
	args := os.Args[2:]

	// Read from stdin to get the additional arguments.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		arg := scanner.Text()
		args = append(args, arg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	// Execute the command with the provided arguments.
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}