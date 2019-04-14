package main_test

import (
	"testing"

	main "github.com/lvnacapital/algorand"
)

func TestMain(t *testing.T) {
	if err := main.Run(); err != nil {
		t.Errorf("Failed to run.")
	}
}
