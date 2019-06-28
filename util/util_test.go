package util_test

import (
	"fmt"
	"testing"

	"github.com/lvnacapital/algorand-go/util"
	"github.com/spf13/viper"
)

func TestIsValidGolden(t *testing.T) {
	t.Parallel()
	golden := []string{
		"7777777777777777777777777777777777777777777777777774MSJUVU",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ",
	}

	for _, addr := range golden {
		if util.IsValidAddress(addr) != true {
			t.Errorf("Expected address %s to be invalid.", golden)
		}
	}
}

func TestMakeClients(t *testing.T) {
	t.Parallel()
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

func TestCodeValid(t *testing.T) {
	t.Parallel()
	codes := []string{"G2AWNZ377XCR44DV", "NBDKVIJVCQSMYICF"}

	for _, code := range codes {
		if err := util.CheckCode(code); err != nil {
			t.Errorf("Code check failed - %v", err)
		}
	}
}

func TestCodeGen(t *testing.T) {
	t.Parallel()
	util.ForceRand = []byte{104, 70, 170, 161, 53, 20, 36, 204}
	expected := "NBDKVIJVCQSMYICF"
	res := util.GenerateCode()
	if res != expected {
		t.Errorf("Code generation failed: %s", res)
	}
}
