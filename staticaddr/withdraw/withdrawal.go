package withdraw

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Withdrawal struct {
	TxID               chainhash.Hash
	DepositOutpoints   []string
	TotalDepositAmount btcutil.Amount
	WithdrawnAmount    btcutil.Amount
	ChangeAmount       btcutil.Amount
	ConfirmationHeight uint32
}
