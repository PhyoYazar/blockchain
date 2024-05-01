package state

import "github.com/PhyoYazar/blockchain/foundation/blockchain/database"

// RetrieveMempool returns a copy of the mempool.
func (s *State) RetrieveMempool() []database.BlockTx {
	return s.mempool.PickBest()
}
