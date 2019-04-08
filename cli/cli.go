package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadLine takes in user input.
func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	resp, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read stdin.")
		os.Exit(1)
	}
	fmt.Printf("\n")
	return strings.TrimSpace(resp)
}

// ClearScreen clears the terminal window and scrollback buffer.
func ClearScreen() {
	// Standard clear command.
	fmt.Printf("\033[H\033[2J")

	// Clear scrollback buffer, if supported.
	fmt.Printf("\033[3J")
}
