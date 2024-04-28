// Package public maintains the group of handlers for public access.
package public

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/PhyoYazar/blockchain/business/web/v1"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/database"
	"github.com/PhyoYazar/blockchain/foundation/blockchain/state"
	"github.com/PhyoYazar/blockchain/foundation/nameservice"
	"github.com/PhyoYazar/blockchain/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of bar ledger endpoints.
type Handlers struct {
	Log   *zap.SugaredLogger
	State *state.State
	NS    *nameservice.NameService
}

// // Sample just provides a starting point for the class.
// func (h Handlers) Sample(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	resp := struct {
// 		Status string
// 	}{
// 		Status: "OK",
// 	}

// 	return web.Respond(ctx, w, resp, http.StatusOK)
// }

// SubmitWalletTransaction adds new transactions to the mempool.
func (h Handlers) SubmitWalletTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	// Decode the JSON in the post call into a Signed transaction.
	var signedTx database.SignedTx
	if err := web.Decode(r, &signedTx); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	h.Log.Infow("add tran", "traceid", v.TraceID, "sig:nonce", signedTx, "from", signedTx.FromID, "to", signedTx.ToID, "value", signedTx.Value, "tip", signedTx.Tip)

	// Ask the state package to add this transaction to the mempool. Only the
	// checks are the transaction signature and the recipient account format.
	// It's up to the wallet to make sure the account has a proper balance and
	// nonce. Fees will be taken if this transaction is mined into a block.
	if err := h.State.UpsertWalletTransaction(signedTx); err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "transactions added to mempool",
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}
