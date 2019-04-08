package main

import (
	"fmt"
	"os"

	"github.com/lvnacapital/algorand/cmd"
)

func main() {
	if err := cmd.AlgorandCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
