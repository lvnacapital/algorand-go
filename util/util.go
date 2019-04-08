package util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
)

const investorCodeLen = 16
const investorCodeLenDecoded = 10
const registrationBaseURL = "https://investorkey.algodev.network/register"

var numKeys int
var genInvestorCode bool
var createWallet bool

func base32Encode(b []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

func base32Decode(b string) ([]byte, error) {
	return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(b)
}

func checkInvestorCode(investorCode string) error {
	if len(investorCode) != investorCodeLen {
		return fmt.Errorf("investor code should be %d characters long", investorCodeLen)
	}

	// Decode investor code
	codeBytes, err := base32Decode(investorCode)
	if err != nil {
		return fmt.Errorf("couldn't parse investor code")
	}

	if len(codeBytes) != investorCodeLenDecoded {
		return fmt.Errorf("invalid investor coded decoded length")
	}

	// Pull out 2-byte checksum
	cksum := codeBytes[investorCodeLenDecoded-2:]

	// Compute expected checksum
	hash := sha512.Sum512_256(codeBytes[:investorCodeLenDecoded-2])

	// Check checksum
	if !bytes.Equal(cksum, hash[:2]) {
		return fmt.Errorf("invalid checksum")
	}

	// Valid investor code
	return nil
}

func generateInvestorCode() (res string) {
	// Generate raw code bytes
	var raw [investorCodeLenDecoded - 2]byte
	_, err := rand.Read(raw[:])
	if err != nil {
		panic(fmt.Sprintf("broken system randomness: %s", err))
	}

	// Compute checksum
	hash := sha512.Sum512_256(raw[:])
	cksum := hash[:2]

	// Append checksum
	rawWithCksum := append(raw[:], cksum...)
	res = base32Encode(rawWithCksum)

	// Check that the generated code is valid
	if len(res) != investorCodeLen {
		panic(fmt.Sprintf("generated bad investor code: %s", res))
	}

	return
}
