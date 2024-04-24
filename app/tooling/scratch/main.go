package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// Need to load the private key file for the configured beneficiary so the
	// account can get credited with fees and tips.
	path := fmt.Sprintf("%s%s.ecdsa", "zblock/accounts/", "kennedy")
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	addr := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println(addr)

	v := struct {
		Name string `json:"name"`
	}{
		Name: "Bill",
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Hash the transaction data into a 32 byte array. This will provide
	// a data length consistency with all transactions.
	txHash := crypto.Keccak256(data)

	// Sign the hash with the private key to produce a signature.
	sig, err := crypto.Sign(txHash, privateKey)
	if err != nil {
		return fmt.Errorf("signature error: %w", err)
	}

	fmt.Println("SIG: ", sig)

	// =========================================================================

	sigPublicKey, err := crypto.Ecrecover(txHash, sig)
	if err != nil {
		return err
	}

	// Capture the public key associated with this signature.
	x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	publicKey := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	// Extract the account address from the public key.

	address := crypto.PubkeyToAddress(publicKey).String()
	fmt.Println(address)

	return nil
}
