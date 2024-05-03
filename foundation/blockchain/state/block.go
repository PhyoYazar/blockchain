package state

import (
	"context"
	"errors"

	"github.com/PhyoYazar/blockchain/foundation/blockchain/database"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/signature"
)

// ErrNoTransactions is returned when a block is requested to be created
// and there are not enough transactions.
var ErrNoTransactions = errors.New("no transactions in mempool")

// =============================================================================

// MineNewBlock attempts to create a new block with a proper hash that can become
// the next block in the chain.
func (s *State) MineNewBlock(ctx context.Context) (database.Block, error) {
	defer s.evHandler("viewer: MineNewBlock: MINING: completed")

	s.evHandler("state: MineNewBlock: MINING: check mempool count")

	// Are there enough transactions in the pool.
	if s.mempool.Count() == 0 {
		return database.Block{}, ErrNoTransactions
	}

	// Pick the best transactions from the mempool.
	trans := s.mempool.PickBest(s.genesis.TransPerBlock)

	// // If PoA is being used, drop the difficulty down to 1 to speed up
	// // the mining operation.
	// difficulty := s.genesis.Difficulty
	// if s.Consensus() == ConsensusPOA {
	// 	difficulty = 1
	// }

	// Attempt to create a new block by solving the POW puzzle. This can be cancelled.
	block, err := database.POW(ctx, database.POWArgs{
		BeneficiaryID: s.beneficiaryID,
		Difficulty:    s.genesis.Difficulty,
		MiningReward:  s.genesis.MiningReward,
		PrevBlock:     s.db.LatestBlock(),
		StateRoot:     signature.ZeroHash,
		Trans:         trans,
		EvHandler:     s.evHandler,
	})
	if err != nil {
		return database.Block{}, err
	}

	// Just check one more time we were not cancelled.
	if ctx.Err() != nil {
		return database.Block{}, ctx.Err()
	}

	s.evHandler("state: MineNewBlock: MINING: validate and update database")

	// // Validate the block and then update the blockchain database.
	// if err := s.validateUpdateDatabase(block); err != nil {
	// 	return database.Block{}, err
	// }

	// HACK THIS FOR NOW
	for _, tx := range trans {
		s.evHandler("state: Remove Tx[%s]", tx)
		s.mempool.Delete(tx)
	}

	return block, nil
}
