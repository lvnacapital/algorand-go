package cmd

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"golang.org/x/crypto/ssh/terminal"
)

// GetWallet finds and returns the wallet handle
func GetWallet() (string, error) {
	// Get the list of wallets
	walletsList, err := kmdClient.ListWallets()
	if err != nil {
		return "", fmt.Errorf("Error listing wallets - %s", err)
	} else if len(walletsList.Wallets) <= 0 {
		return "", fmt.Errorf("No wallets available")
	}

	var walletID string
	fmt.Printf("\nHave %d wallet(s):\n", len(walletsList.Wallets))
	for i, wallet := range walletsList.Wallets {
		if WalletName != "" { // Find our wallet name in the list
			if wallet.Name == WalletName {
				fmt.Printf("Found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
				walletID = wallet.ID
				break
			}
		} else { // List wallets for selection
			fmt.Printf("[%d] Name: %s\tID: %s\n", i+1, wallet.Name, wallet.ID)
		}
	}
	if walletID == "" {
		term := terminal.NewTerminal(os.Stdin, "")
		for {
			if len(walletsList.Wallets) == 1 {
				fmt.Printf("Select wallet [%s]: ", "1")
			} else {
				fmt.Printf("Select wallet [%s%d]: ", "1-", len(walletsList.Wallets))
			}
			walletNum, err := term.ReadLine()
			if err != nil {
				return "", fmt.Errorf("Error getting wallet number: %s", err)
			}
			i, err := strconv.Atoi(string(walletNum))
			if err != nil || i > len(walletsList.Wallets) || i <= 0 {
				fmt.Print("Invalid wallet number. Please try again.\n")
				continue
			}
			walletID = walletsList.Wallets[i-1].ID
			WalletName = walletsList.Wallets[i-1].Name
			break
		}
	}
	// fmt.Printf("Picked wallet %s.\n", walletID)

	if WalletPassword == "" {
		fmt.Printf("Please type in the password for '%s': ", WalletName)
		pw, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", fmt.Errorf("\nError getting password: %s", err)
		}
		WalletPassword = string(pw)
		fmt.Print("\n")
	}

	// Get a wallet handle
	initRes, err := kmdClient.InitWalletHandle(walletID, WalletPassword)
	if err != nil {
		return "", fmt.Errorf("\nError initializing wallet handle: %s", err)
	}
	walletHandle := initRes.WalletHandleToken

	return walletHandle, nil
}

func getPrivateKey() (keyBytes []byte, err error) {
	if WalletMnemonic == "" {
		term := terminal.NewTerminal(os.Stdin, "")
		for {
			fmt.Print("\nEnter the wallet mnemonic: ")
			m, err := term.ReadLine()
			if err != nil {
				return nil, fmt.Errorf("Error getting mnemonic - %v", err)
			}
			WalletMnemonic = string(m)
			if keyBytes, err = mnemonic.ToKey(WalletMnemonic); err != nil {
				fmt.Printf("Failed to get key. Try again - %v", err)
				continue
			}
			break
		}
	} else {
		if keyBytes, err = mnemonic.ToKey(WalletMnemonic); err != nil {
			return nil, fmt.Errorf("Failed to get key from -m: %v", err)
		}
	}

	return
}
