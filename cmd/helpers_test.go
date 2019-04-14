package cmd_test

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var mnemonic = "fire enlist diesel stamp nuclear chunk student stumble call snow flock brush example slab guide choice option recall south kangaroo hundred matrix school above zero"
var walletName = "testwallet"
var walletPassword = "testpassword"
var algodAvailable = false
var kmdAvailable = true

func init() {
	if os.Getenv("CI") == "true" {
		algodAvailable = httpAlgod()
		kmdAvailable = httpKmd()
	}

	// walletHandle, err := cmd.getWallet()
	// if err != nil {
	// 	return
	// }
}

// var keyBytes = []byte{185, 178, 18, 123, 70, 173, 203, 162, 236, 186, 215, 3, 97, 179, 178, 210, 49, 167, 43, 243, 44, 40, 221, 220, 108, 160, 153, 151, 183, 36, 26, 152}

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

func httpAlgod() bool {
	url := fmt.Sprintf("http://%s:%s/v1/status", os.Getenv("ALGORAND_HOST"), os.Getenv("ALGOD_PORT"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Algo-API-Token", os.Getenv("ALGOD_TOKEN"))
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return false
	}
	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))

	return true
}

func httpKmd() bool {
	url := fmt.Sprintf("http://%s:%s/v1/wallets", os.Getenv("ALGORAND_HOST"), os.Getenv("KMD_PORT"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-KMD-API-Token", os.Getenv("KMD_TOKEN"))
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return false
	}
	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))

	return true
}
