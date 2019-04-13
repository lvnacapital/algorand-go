package util_test

import (
	"fmt"
	"testing"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/viper"
)

func TestIsValidGolden(t *testing.T) {
	golden := "7777777777777777777777777777777777777777777777777774MSJUVU"

	if util.IsValidAddress(golden) != true {
		t.Errorf("Expected address %s to be invalid.", golden)
	}
}

func TestMakeClients(t *testing.T) {
	nodeConfig := util.Node{
		AlgodAddress: fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port")),
		KmdAddress:   fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port")),
		AlgodToken:   viper.GetString("algod-token"),
		KmdToken:     viper.GetString("kmd-token"),
	}

	if _, _, err := util.MakeClients(&nodeConfig); err != nil {
		t.Errorf("Failed to make clients.")
	}
}
