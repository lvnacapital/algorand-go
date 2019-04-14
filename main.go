package main

import (
	"fmt"
	"os"

	"github.com/lvnacapital/algorand/cmd"
)

// Run executes the top-level command
func Run() error {
	if err := cmd.AlgorandCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func main() {
	Run()
}
