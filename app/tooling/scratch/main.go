package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	err := sign()

	if err != nil {
		log.Fatalln(err)
	}
}

func sign() error {
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

	data, err := stamp(v)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}

	// Sign the hash with the private key to produce a signature.
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return fmt.Errorf("signature error: %w", err)
	}

	fmt.Printf("SIG: 0x%s\n", hex.EncodeToString(sig))

	// =========================================================================

	v2 := struct {
		Name string `json:"name"`
	}{
		Name: "Bills",
	}

	data2, err := stamp(v2)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}

	// =========================================================================

	sigPublicKey, err := crypto.Ecrecover(data2, sig)
	if err != nil {
		return err
	}

	// Check the public key validates the data and signature.
	rs := sig[:crypto.RecoveryIDOffset]
	if !crypto.VerifySignature(sigPublicKey, data2, rs) {
		return errors.New("invalid signature produced")
	}

	// Capture the public key associated with this signature.
	x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	publicKey := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	// Extract the account address from the public key.
	address := crypto.PubkeyToAddress(publicKey).String()
	fmt.Println(address)

	return nil
}

// =============================================================================

// stamp returns a hash of 32 bytes that represents this data with
// the Ardan stamp embedded into the final hash.
func stamp(value any) ([]byte, error) {

	// Marshal the data.
	v, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// This stamp is used so signatures we produce when signing data
	// are always unique to the Ardan blockchain.
	stamp := []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(v)))

	// Hash the stamp and txHash together in a final 32 byte array
	// that represents the data.
	data := crypto.Keccak256(stamp, v)

	return data, nil
}
