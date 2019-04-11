package util_test

import (
	"testing"

	"github.com/lvnacapital/algorand/util"
)

func TestIsValidGolden(t *testing.T) {
	golden := "7777777777777777777777777777777777777777777777777774MSJUVU"

	if util.IsValidAddress(golden) != true {
		t.Errorf("Expected address %s to be invalid.", golden)
	}
}
