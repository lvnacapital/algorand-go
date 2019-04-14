package cmd_test

import (
	"bytes"

	"github.com/spf13/cobra"
)

var mnemonic = "fire enlist diesel stamp nuclear chunk student stumble call snow flock brush example slab guide choice option recall south kangaroo hundred matrix school above zero"
var walletName = "testwallet"
var walletPassword = "testpassword"

// var keyBytes []byte = [185 178 18 123 70 173 203 162 236 186 215 3 97 179 178 210 49 167 43 243 44 40 221 220 108 160 153 151 183 36 26 152]

// func init() {
// 	_, filename, _, _ := runtime.Caller(0)
// 	dir := path.Join(path.Dir(filename), "..")
// 	err := os.Chdir(dir)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func emptyRun(*cobra.Command, []string) {}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
