// Package database handles all the lower level support for maintaining the
// blockchain in storage and maintaining an in-memory databse of account information.
package database

import (
	"sync"

	"github.com/PhyoYazar/blockchain/foundation/blockchain/genesis"
)

// Storage interface represents the behavior required to be implemented by any
// package providing support for reading and writing the blockchain.
type Storage interface {
	Write(blockData BlockData) error
	GetBlock(num uint64) (BlockData, error)
	Close() error
}

// =============================================================================

// Database manages data related to accounts who have transacted on the blockchain.
type Database struct {
	mu      sync.RWMutex
	genesis genesis.Genesis
	// latestBlock Block
	accounts map[AccountID]Account
	storage  Storage
}

// New constructs a new database and applies account genesis information and
// reads/writes the blockchain database on disk if a dbPath is provided.
func New(genesis genesis.Genesis, storage Storage, evHandler func(v string, args ...any)) (*Database, error) {
	db := Database{
		genesis:  genesis,
		accounts: make(map[AccountID]Account),
		storage:  storage,
	}

	// Update the database with account balance information from genesis.
	for accountStr, balance := range genesis.Balances {
		accountID, err := ToAccountID(accountStr)
		if err != nil {
			return nil, err
		}
		db.accounts[accountID] = newAccount(accountID, balance)
	}

	return &db, nil
}

// Close closes the open blocks database.
func (db *Database) Close() {
	db.storage.Close()
}
