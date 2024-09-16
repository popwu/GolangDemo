package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sour-is/bitcoin/address"
	"github.com/sour-is/bitcoin/bip38"
)

func main() {
	info()
	fmt.Println("================================")
	testDecrypted()
}

func testDecrypted() {
	decryptedKey := "6PRWa5APpFFAczW8nwTMhPBLF5We2yKRDwMTs4q9m6uCrPFgEN7dE8TFnK"
	passphrase := "TestingOneTwoThree"
	private, err := bip38.Decrypt(decryptedKey, passphrase)
	if err != nil {
		panic(err)
	}
	fmt.Println("Private Key Bytes (Hex)", hex.EncodeToString(private.Bytes()))
}

func info() {
	privateHex := "c473061e73df28bd13cbaa5deeb0c6a6ed83ed6d449048905741586981c3ec50"
	fmt.Println("Length of privateHex:", len(privateHex))

	private, err := address.ReadPrivateKey(privateHex)
	if err != nil {
		panic(err)
	}
	fmt.Println("Private Key:", private)
	fmt.Println("Private Key Bytes (Hex):", hex.EncodeToString(private.Bytes()))
	fmt.Println("Private Key Address:", private.Address())

	// Encrypt a private key
	encryptedKey := bip38.Encrypt(private, "TestingOneTwoThree")
	fmt.Println("Encrypted Key:", encryptedKey)

	// Decrypt the encrypted key
	decryptedKey, err := bip38.Decrypt(encryptedKey, "TestingOneTwoThree")
	if err != nil {
		panic(err)
	}
	fmt.Println("Decrypted Key:", decryptedKey)
}
