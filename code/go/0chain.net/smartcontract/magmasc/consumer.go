package magmasc

import (
	"encoding/json"
	"fmt"

	"0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/util"
)

// registerConsumer represents registerConsumer MagmaSmartContract function and allows registering Consumer in blockchain.
//
// registerConsumer creates Consumer with Consumer.ID (equals to transaction client ID),
// adds it to all Consumers list, creates stakePool for new Consumer and saves results in provided state.StateContextI.
func (msc *MagmaSmartContract) registerConsumer(txn *transaction.Transaction, balances state.StateContextI) (string, error) {
	const errCode = "add_consumer"

	consumers, err := extractConsumers(balances)
	if err != nil {
		return "", common.NewErrorf(errCode, "retrieving all consumers from state failed with error: %v ", err)
	}

	var (
		consumer = Consumer{
			ID: txn.ClientID,
		}
	)
	if containsConsumer(msc.ID, consumer, consumers, balances) {
		return "", common.NewErrorf(errCode, "consumer with id=`%s` already exist", consumer.ID)
	}

	if err := createAndInsertConsumerStakePool(consumer.ID, msc.ID, balances); err != nil {
		return "", common.NewErrorf(errCode, "creating stake pool for consumer failed with err: %v", err)
	}

	// save the all consumers
	consumers.Nodes.add(&consumer)
	_, err = balances.InsertTrieNode(AllConsumersKey, consumers)
	if err != nil {
		return "", common.NewErrorf(errCode, "saving the all consumers failed with error: %v ", err)
	}

	// save the new consumer
	_, err = balances.InsertTrieNode(consumer.GetKey(msc.ID), &consumer)
	if err != nil {
		return "", common.NewErrorf(errCode, "saving consumer failed with error: %v ", err)
	}

	return "", nil
}

// extractConsumers extracts all Consumers represented in JSON bytes stored in state.StateContextI with AllConsumersKey.
//
// extractConsumers returns err if state.StateContextI does not contain Consumers or stored Consumers bytes have invalid
// format.
func extractConsumers(balances state.StateContextI) (*Consumers, error) {
	consumers := &Consumers{}
	consumersBytes, err := balances.GetTrieNode(AllConsumersKey)
	if consumersBytes == nil || err != nil {
		return consumers, err
	}

	err = json.Unmarshal(consumersBytes.Encode(), consumers)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}
	return consumers, nil
}

// containsConsumer looks for provided Consumer in provided Consumers and state.StateContextI.
// If Consumer will be found it returns true, else false.
func containsConsumer(scKey string, consumer Consumer, consumers *Consumers, balances state.StateContextI) bool {
	for _, b := range consumers.Nodes {
		if b.ID == consumer.ID || b.BaseURL == consumer.BaseURL {
			return true
		}
	}

	_, err := balances.GetTrieNode(consumer.GetKey(scKey))
	if err == nil {
		return true
	}

	return false
}

// createAndInsertConsumerStakePool creates stakePool for Consumer and saves it in state.StateContextI.
//
// if stakePool for provided Consumer.ID already exist it returns ErrStakePoolExist. Also, createAndInsertConsumerStakePool
// returns err occurred while inserting new stakePool in state.StateContextI.
func createAndInsertConsumerStakePool(consumerID, scKey string, balances state.StateContextI) error {
	_, err := balances.GetTrieNode(stakePoolKey(scKey, consumerID))
	if err != util.ErrValueNotPresent {
		return ErrStakePoolExist
	}

	sp := new(stakePool)
	sp.ID = consumerID

	_, err = balances.InsertTrieNode(stakePoolKey(scKey, consumerID), sp)
	if err != nil {
		return err
	}

	return nil
}
