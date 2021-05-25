package magmasc

import (
	"0chain.net/chaincore/tokenpool"
	"0chain.net/core/datastore"
)

// stakePool represents simple stake pool implementation used in MagmaSmartContract.
type stakePool struct {
	*tokenpool.ZcnPool
}

// stakePoolKey represents uniq stakePool key used to saving stakePool in state.StateContextI.
//
// Resulting key represents concatenated smart contract key, ":stakepool:" and node id.
func stakePoolKey(scKey, id string) datastore.Key {
	return scKey + ":stakepool:" + id
}
