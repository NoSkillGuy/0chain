package magmasc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/util"
	"0chain.net/smartcontract"
)

// registerProvider represents registerProvider MagmaSmartContract function and allows registering Provider in blockchain.
//
// registerProvider creates Provider with Provider.ID (equals to transaction client GetID),
// adds it to all Nodes list and saves results in provided state.StateContextI.
func (msc *MagmaSmartContract) registerProvider(txn *transaction.Transaction,
	input []byte, balances state.StateContextI) (string, error) {
	const errCode = "register_provider"

	providers, err := extractProviders(balances)
	if err != nil {
		return "", common.NewErrorf(errCode, "retrieving all providers from state failed with error: %v ", err)
	}

	provider := Provider{}
	if err := json.Unmarshal(input, &provider); err != nil {
		return "", common.NewErrorf(errCode, "unmarshalling input failed with err: %v", err)
	}
	provider.ID = txn.ClientID
	if containsProvider(msc.ID, provider, providers, balances) {
		return "", common.NewErrorf(errCode, "provider with id=`%s` already exist", provider.ID)
	}

	// save the all providers
	providers.Nodes.add(&provider)
	_, err = balances.InsertTrieNode(AllProvidersKey, providers)
	if err != nil {
		return "", common.NewErrorf(errCode, "saving the all providers failed with error: %v ", err)
	}

	// save the new provider
	_, err = balances.InsertTrieNode(nodeKey(msc.ID, provider.ID, providerType), &provider)
	if err != nil {
		return "", common.NewErrorf(errCode, "saving provider failed with error: %v ", err)
	}

	return string(provider.Encode()), nil
}

// extractProviders extracts all provider Nodes represented in JSON bytes stored in state.StateContextI with AllProvidersKey.
//
// extractProviders returns err if state.StateContextI does not contain Nodes or stored Nodes bytes have invalid
// format.
func extractProviders(balances state.StateContextI) (*Providers, error) {
	providers := &Providers{}
	providerNV, err := balances.GetTrieNode(AllProvidersKey)
	if err != nil && err != util.ErrValueNotPresent {
		return nil, err
	}
	if err == util.ErrValueNotPresent || providerNV == nil {
		return providers, nil
	}

	if err := json.Unmarshal(providerNV.Encode(), providers); err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}

	return providers, nil
}

// getProviderTerms represents MagmaSmartContract handler. getProviderTerms looks for Provider with id, passed in params url.Values,
// in provided state.StateContextI and returns Provider.Terms.
func (msc *MagmaSmartContract) getProviderTerms(_ context.Context, params url.Values, balances state.StateContextI) (interface{}, error) {
	id := params.Get("provider_id")

	provider, err := extractProvider(id, msc.ID, balances)
	if err != nil {
		return nil, smartcontract.NewErrNoResourceOrErrInternal(err, true,
			"extracting provider from state failed with err")
	}

	return provider.Terms, nil
}

// extractProviders extracts Provider represented in JSON bytes stored in state.StateContextI.
//
// extractProviders returns err if state.StateContextI does not contain Nodes or stored Nodes bytes have invalid format.
func extractProvider(id, scKey string, balances state.StateContextI) (*Provider, error) {
	providerNV, err := balances.GetTrieNode(nodeKey(scKey, id, providerType))
	if err != nil {
		return nil, err
	}

	provider := new(Provider)
	if err := json.Unmarshal(providerNV.Encode(), provider); err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}

	return provider, nil
}
