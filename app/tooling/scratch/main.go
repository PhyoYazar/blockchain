package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PhyoYazar/blockchain/foundation/blockchain/database"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/merkle"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/signature"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/storage/disk"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// err := sign()
	// err := writeBlock()
	err := readBlock()

	if err != nil {
		log.Fatalln(err)
	}
}

func readBlock() error {
	d, err := disk.New("zblock/miner1")
	if err != nil {
		return nil
	}

	blockData, err := d.GetBlock(1)
	if err != nil {
		return nil
	}

	fmt.Println(blockData)

	block, err := database.ToBlock(blockData)
	if err != nil {
		return nil
	}

	if blockData.Header.TransRoot != block.MerkleTree.RootHex() {
		return errors.New("merkle tree comparison failed")
	}

	fmt.Println("markle tree is matched!")

	return nil
}

func writeBlock() error {
	txs := []database.Tx{
		{
			ChainID: 1,
			Nonce:   1,
			FromID:  "0xdd6B972ffcc631a62CAE1BB9d80b7ff429c8ebA4",
			ToID:    "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
			Value:   100,
			Tip:     50,
		},
		{
			ChainID: 1,
			Nonce:   2,
			FromID:  "0xdd6B972ffcc631a62CAE1BB9d80b7ff429c8ebA4",
			ToID:    "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
			Value:   100,
			Tip:     50,
		},
	}

	blockTxs := make([]database.BlockTx, len(txs))

	for i, tx := range txs {
		blockTx, err := signToBlockTx(tx, 15)
		if err != nil {
			return err
		}

		blockTxs[i] = blockTx
	}

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(blockTxs)
	if err != nil {
		return err
	}

	beneficiaryID, err := database.ToAccountID("0xF01813E4B85e178A83e29B8E7bF26BD830a25f32")
	if err != nil {
		return err
	}

	// Construct the block to be mined.
	block := database.Block{
		Header: database.BlockHeader{
			Number:        1,
			PrevBlockHash: signature.ZeroHash,
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			BeneficiaryID: beneficiaryID,
			Difficulty:    6,
			MiningReward:  700,
			StateRoot:     "not defined",
			TransRoot:     tree.RootHex(), //
			Nonce:         0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	bd := database.NewBlockData(block)

	d, err := disk.New("zblock/miner1")
	if err != nil {
		return nil
	}

	if err := d.Write(bd); err != nil {
		return nil
	}

	return nil
}

func signToBlockTx(tx database.Tx, gas uint64) (database.BlockTx, error) {
	pk, err := crypto.HexToECDSA("fae85851bdf5c9f49923722ce38f3c1defcfd3619ef5453230a58ad805499959")
	if err != nil {
		return database.BlockTx{}, err
	}

	signedTx, err := tx.Sign(pk)
	if err != nil {
		return database.BlockTx{}, err
	}

	return database.NewBlockTx(signedTx, gas, 1), nil
}

// =============================================================================

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
