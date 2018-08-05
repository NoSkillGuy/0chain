package chain

import (
	"bytes"
	"context"

	"0chain.net/block"
	"0chain.net/common"
	"0chain.net/datastore"
	. "0chain.net/logging"
	"0chain.net/state"
	"0chain.net/transaction"
	"0chain.net/util"
	"go.uber.org/zap"
)

//StateMismatch - indicate if there is a mismatch between computed state and received state of a block
const StateMismatch = "state_mismatch"

/*ComputeState - compute the state for the block */
func (c *Chain) ComputeState(ctx context.Context, b *block.Block) error {
	if b.IsStateComputed() {
		return nil
	}
	/*TODO - this needs to be robust. Will lead to stackoverflow at the moment as we are not strict on state correctness
	if b.PrevBlock != nil {
		pbState := b.PrevBlock.GetBlockState()
		if !b.PrevBlock.IsStateComputed() {
			Logger.Info("compute state - previous block state not ready", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.String("prev_block", b.PrevHash), zap.Int8("prev_block_state", pbState))
			c.ComputeState(ctx, b.PrevBlock)
		}
	} else {
		Logger.Error("compute state - previous block not available", zap.Int64("round", b.Round), zap.String("block", b.Hash))
		return ErrPreviousBlockUnavailable
	}*/
	c.rebaseState()
	for _, txn := range b.Txns {
		if datastore.IsEmpty(txn.ClientID) {
			txn.ComputeClientID()
		}
		if !c.UpdateState(b, txn) {
			return common.NewError("state_update_error", "error updating state")
		}
	}
	if bytes.Compare(b.ClientStateHash, b.ClientState.GetRoot()) != 0 {
		Logger.Error("validate transaction state hash error", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Int("block_size", len(b.Txns)), zap.Int("changes", len(b.ClientState.GetChangeCollector().GetChanges())), zap.String("block_state_hash", util.ToHex(b.ClientStateHash)), zap.String("computed_state_hash", util.ToHex(b.ClientState.GetRoot())))
		return common.NewError(StateMismatch, "computed state hash doesn't match with the state hash of the block")
	}
	b.SetStateIsComputed(true)
	return nil
}

func (c *Chain) rebaseState() {
	lfb := c.LatestFinalizedBlock
	ndb := lfb.ClientState.GetNodeDB()
	if ndb != c.StateDB {
		lfb.ClientState.SetNodeDB(c.StateDB)
		if lndb, ok := ndb.(*util.LevelNodeDB); ok {
			Logger.Debug("finalize round - rebasing current state db", zap.Int64("round", lfb.Round), zap.String("block", lfb.Hash), zap.String("hash", util.ToHex(lfb.ClientState.GetRoot())))
			lndb.RebaseCurrentDB(c.StateDB)
			Logger.Debug("finalize round - rebased current state db", zap.Int64("round", lfb.Round), zap.String("block", lfb.Hash), zap.String("hash", util.ToHex(lfb.ClientState.GetRoot())))
		}
	}
}

/*UpdateState - update the state of the transaction w.r.t the given block
* The block starts off with the state from the prior block and as transactions are processed into a block, the state gets updated
* If a state can't be updated (e.g low balance), then a false is returned so that the transaction will not make it into the block
 */
func (c *Chain) UpdateState(b *block.Block, txn *transaction.Transaction) bool {
	clientState := b.ClientState
	fs, err := c.getState(clientState, txn.ClientID)
	if err != nil {
		if err != util.ErrValueNotPresent {
			Logger.Debug("update state", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Int8("block_state", b.GetBlockState()), zap.Any("txn", txn), zap.Error(err))
			return false
		} else {
			Logger.Debug("update state (value not present)", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Int8("block_state", b.GetBlockState()), zap.Any("txn", txn), zap.Error(err))
		}
	}
	tbalance := state.Balance(txn.Value)
	switch txn.TransactionType {
	case transaction.TxnTypeSend:
		if fs.Balance < tbalance {
			Logger.Debug("low balance", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Any("state", fs), zap.Any("txn", txn))
			return false
		}
		fs.Balance -= tbalance
		if fs.Balance == 0 {
			_, err = clientState.Delete(util.Path(txn.ClientID))
		} else {
			_, err = clientState.Insert(util.Path(txn.ClientID), fs)
		}
		if err != nil {
			Logger.Debug("update state - error", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Any("txn", datastore.ToJSON(txn)), zap.Error(err))
		}
		ts, err := c.getState(clientState, txn.ToClientID)
		if err != nil {
			Logger.Debug("update state (to client)", zap.Any("txn", datastore.ToJSON(txn)), zap.Error(err))
			return false
		}
		ts.Balance += tbalance
		clientState.Insert(util.Path(txn.ToClientID), ts)
		return true
	default:
		return true // TODO: This should eventually return false by default for all unkown cases
	}
}

func (c *Chain) getState(clientState util.MerklePatriciaTrieI, clientID string) (*state.State, error) {
	if clientState == nil {
		return nil, common.NewError("get state", "client state does not exist")
	}
	s := &state.State{}
	s.Balance = state.Balance(0)
	ss, err := clientState.GetNodeValue(util.Path(clientID))
	if err != nil {
		if err != util.ErrValueNotPresent {
			return nil, err
		}
	} else {
		s = c.ClientStateDeserializer.Deserialize(ss).(*state.State)
	}
	return s, nil
}
