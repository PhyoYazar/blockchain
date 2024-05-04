package state

import (
	"github.com/PhyoYazar/blockchain/foundation/blockchain/database"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/peer"
)

// RetrieveHost returns a copy of host information.
func (s *State) RetrieveHost() string {
	return s.host
}

// RetrieveMempool returns a copy of the mempool.
func (s *State) RetrieveMempool() []database.BlockTx {
	return s.mempool.PickBest()
}

// RetrieveLatestBlock returns a copy the current latest block.
func (s *State) RetrieveLatestBlock() database.Block {
	return s.db.LatestBlock()
}

// RetrieveAccounts returns a copy of the database accounts.
func (s *State) RetrieveAccounts() map[database.AccountID]database.Account {
	return s.db.CopyAccounts()
}

// RetrieveKnownPeers retrieves a copy of the known peer list.
func (s *State) RetrieveKnownPeers() []peer.Peer {
	return s.knownPeers.Copy(s.host)
}
