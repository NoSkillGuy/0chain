package magmasc

import (
	"encoding/json"

	"0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
)

// acceptTerms represent MagmaSmartContract function. acceptTerms checks input for validity,
// sets the client's id from transaction to Acknowledgment.ConsumerID, set's hash of transaction to Acknowledgment.ID
// and inserts resulted Acknowledgment in provided state.StateContextI.
func (msc *MagmaSmartContract) acceptTerms(txn *transaction.Transaction, input []byte, balances state.StateContextI) (string, error) {
	const errCode = "accept_terms"

	acknol := new(Acknowledgment)
	if err := json.Unmarshal(input, acknol); err != nil {
		return "", common.NewErrorf(errCode, "unmarshalling input failed with err: %v", err)
	}
	if !acknol.isValid() {
		return "", common.NewErrorf(errCode, "provided acknowledgment is invalid")
	}
	acknol.ConsumerID = txn.ClientID
	acknol.ID = txn.GetKey()

	if _, err := balances.InsertTrieNode(acknowledgmentKey(msc.ID, acknol.ID), acknol); err != nil {
		return "", common.NewErrorf(errCode, "saving acknowledgment failed with error: %v ", err)
	}

	return string(acknol.Encode()), nil
}
