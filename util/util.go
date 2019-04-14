package util

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
)

const codeLength = 16
const codeLenghtDecoded = 10
const algorandAddressLength = 58

// ForceRand for deterministic test
var ForceRand []byte

// Node configuration
type Node struct {
	AlgodAddress string
	KmdAddress   string
	AlgodToken   string
	KmdToken     string
}

// IsValidAddress ...
func IsValidAddress(address string) bool {
	if reflect.TypeOf(address).String() != "string" {
		return false
	}

	if len(address) != algorandAddressLength {
		return false
	}

	if _, err := types.DecodeAddress(address); err != nil {
		return false
	}

	return true
}

// MakeClients ...
func MakeClients(node *Node) (algodClient algod.Client, kmdClient kmd.Client, err error) {
	// Create an algod client
	if algodClient, err = algod.MakeClient(node.AlgodAddress, node.AlgodToken); err != nil {
		return
	}
	// fmt.Println("Made an algod client")

	// Create a kmd client
	if kmdClient, err = kmd.MakeClient(node.KmdAddress, node.KmdToken); err != nil {
		return
	}
	// fmt.Println("Made a kmd client")

	return
}

// ReadLine takes in user input
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

// ClearScreen clears the terminal window and scrollback buffer
func ClearScreen() {
	if runtime.GOOS != "windows" {
		// Standard clear command.
		fmt.Printf("\033[H\033[2J")

		// Clear scrollback buffer, if supported.
		fmt.Printf("\033[3J")
	}
}

func base32Encode(b []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

func base32Decode(b string) ([]byte, error) {
	return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(b)
}

// CheckCode verifies that the input code and checksum match
func CheckCode(code string) error {
	if len(code) != codeLength {
		return fmt.Errorf("Code should be %d characters long", codeLength)
	}

	// Decode code
	codeBytes, err := base32Decode(code)
	if err != nil {
		return fmt.Errorf("Could not parse code")
	} else if len(codeBytes) != codeLenghtDecoded {
		return fmt.Errorf("Invalid code decoded length")
	}

	// Pull out 2-byte checksum
	cksum := codeBytes[codeLenghtDecoded-2:]

	// Compute expected checksum
	hash := sha512.Sum512_256(codeBytes[:codeLenghtDecoded-2])

	// Check checksum
	if !bytes.Equal(cksum, hash[:2]) {
		return fmt.Errorf("Invalid checksum")
	}

	// Valid code
	return nil
}

// GenerateCode generates a code
func GenerateCode() (res string) {
	// Generate raw code bytes
	raw := make([]byte, codeLenghtDecoded-2)
	if ForceRand != nil {
		raw = ForceRand
	} else {
		if _, err := rand.Read(raw); err != nil {
			panic(fmt.Sprintf("Broken system randomness: %s", err))
		}
	}

	// Compute checksum
	hash := sha512.Sum512_256(raw)
	cksum := hash[:2]

	// Append checksum
	rawWithCksum := append(raw, cksum...)
	res = base32Encode(rawWithCksum)

	// Check that the generated code is valid
	if len(res) != codeLength {
		panic(fmt.Sprintf("Generated bad code: %s", res))
	}

	return
}
