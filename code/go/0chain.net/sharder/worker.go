package sharder

import (
	"0chain.net/chaincore/block"
	"0chain.net/chaincore/chain"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"

	"0chain.net/chaincore/round"
	"0chain.net/core/datastore"
	"0chain.net/core/ememorystore"
	. "0chain.net/core/logging"
	"0chain.net/core/persistencestore"
	"context"
)

/*SetupWorkers - setup the background workers */
func SetupWorkers(ctx context.Context) {
	sc := GetSharderChain()
	go sc.BlockWorker(ctx)              // 1) receives incoming blocks from the network
	go sc.FinalizeRoundWorker(ctx, sc)  // 2) sequentially finalize the rounds
	go sc.FinalizedBlockWorker(ctx, sc) // 3) sequentially processes finalized blocks

	// Setup the deep and proximity scan
	go sc.HealthCheckSetup(ctx, DeepScan)
	go sc.HealthCheckSetup(ctx, ProximityScan)

	go sc.PruneStorageWorker(ctx, time.Minute*5, sc.getPruneCountRoundStorage(), sc.MagicBlockStorage)
	go sc.RegisterSharderKeepWorker(ctx)
}

/*BlockWorker - stores the blocks */
func (sc *Chain) BlockWorker(ctx context.Context) {
	for true {
		select {
		case <-ctx.Done():
			return
		case b := <-sc.GetBlockChannel():
			sc.processBlock(ctx, b)
		}
	}
}

func (sc *Chain) hasRoundSummary(ctx context.Context, rNum int64) (*round.Round, bool) {
	r, err := sc.GetRoundFromStore(ctx, rNum)
	if err == nil {
		return r, true
	}
	return nil, false
}

func (sc *Chain) hasBlockSummary(ctx context.Context, bHash string) (*block.BlockSummary, bool) {
	bSummaryEntityMetadata := datastore.GetEntityMetadata("block_summary")
	bctx := ememorystore.WithEntityConnection(ctx, bSummaryEntityMetadata)
	defer ememorystore.Close(bctx)
	bs, err := sc.GetBlockSummary(bctx, bHash)
	if err == nil {
		return bs, true
	}
	return nil, false
}

func (sc *Chain) hasBlock(bHash string, rNum int64) (*block.Block, bool) {
	b, err := sc.GetBlockFromStore(bHash, rNum)
	if err == nil {
		return b, true
	}
	return nil, false
}

func (sc *Chain) hasBlockTransactions(ctx context.Context, b *block.Block) bool {
	txnSummaryEntityMetadata := datastore.GetEntityMetadata("txn_summary")
	tctx := persistencestore.WithEntityConnection(ctx, txnSummaryEntityMetadata)
	defer persistencestore.Close(tctx)
	for _, txn := range b.Txns {
		_, err := sc.GetTransactionSummary(tctx, txn.Hash)
		if err != nil {
			return false
		}
	}
	return true
}

func (sc *Chain) hasTransactions(ctx context.Context, bs *block.BlockSummary) bool {
	if bs == nil {
		return false
	}
	count, err := sc.getTxnCountForRound(ctx, bs.Round)
	if err != nil {
		return false
	}
	return count == bs.NumTxns
}

func (sc *Chain) RegisterSharderKeepWorker(ctx context.Context) {
	timerCheck := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timerCheck.C:
			if !sc.ActiveInChain() || !sc.IsRegisteredSharderKeep() {
				for !sc.IsRegisteredSharderKeep() {
					txn, err := sc.RegisterSharderKeep()
					if err != nil {
						Logger.Error("register_sharder_keep_worker", zap.Error(err))
					} else {
						if txn == nil || sc.ConfirmTransaction(txn) {
							Logger.Info("register_sharder_keep_worker -- registered")
						} else {
							Logger.Debug("register_sharder_keep_worker -- failed to confirm transaction", zap.Any("txn", txn))
						}
					}
					time.Sleep(time.Second)
				}
			}
		}
	}
}

func (sc *Chain) getPruneCountRoundStorage() func(storage round.RoundStorage) int {
	viper.SetDefault("server_chain.round_magic_block_storage.prune_below_count", chain.DefaultCountPruneRoundStorage)
	pruneBelowCountMB := viper.GetInt("server_chain.round_magic_block_storage.prune_below_count")
	return func(storage round.RoundStorage) int {
		switch storage {
		case sc.MagicBlockStorage:
			return pruneBelowCountMB
		default:
			return chain.DefaultCountPruneRoundStorage
		}
	}
}
