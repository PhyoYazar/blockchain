package state

import (
	"errors"

	"github.com/PhyoYazar/blockchain/foundation/blockchain/database"
)

// =============================================================================

// QueryAccounts returns a copy of the account from the database.
func (s *State) QueryAccounts(account database.AccountID) (database.Account, error) {
	accounts := s.db.CopyAccounts()

	if info, exists := accounts[account]; exists {
		return info, nil
	}

	return database.Account{}, errors.New("not found")
}

// QueryMempoolLength returns the current length of the mempool.
func (s *State) QueryMempoolLength() int {
	return s.mempool.Count()
}
