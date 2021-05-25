package magmasc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rcrowley/go-metrics"

	chainstate "0chain.net/chaincore/chain/state"
	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
)

const (
	Address = "address"

	Name = "magma"
)

// MagmaSmartContract represents smartcontractinterface.SmartContractInterface implementation allows interacting with Magma.
type MagmaSmartContract struct {
	*sci.SmartContract
}

var (
	// Ensure MagmaSmartContract implements smartcontractinterface.SmartContractInterface.
	_ sci.SmartContractInterface = (*MagmaSmartContract)(nil)
)

// NewMagmaSmartContract creates smartcontractinterface.SmartContractInterface implemented by MagmaSmartContract with configured
// RestHandlers and SmartContractExecutionStats.
func NewMagmaSmartContract() sci.SmartContractInterface {
	var sscCopy = &MagmaSmartContract{
		SmartContract: sci.NewSC(Address),
	}
	sscCopy.setSC(sscCopy.SmartContract)
	return sscCopy
}

// These contants represents SmartContractExecutionStats keys, used to identify smart contract functions.
const (
	// registerConsumer represents name for Consumer's registration function.
	registerConsumer = "register_consumer"
)

// setSC sets provided smartcontractinterface.SmartContract to corresponding MagmaSmartContract field
// and configures RestHandlers and SmartContractExecutionStats.
func (msc *MagmaSmartContract) setSC(sc *sci.SmartContract) {
	msc.SmartContract = sc

	// consumer
	msc.SmartContract.SmartContractExecutionStats[registerConsumer] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", msc.ID, registerConsumer), nil)
}

// GetName implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetName() string {
	return Name
}

// GetAddress implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetAddress() string {
	return Address
}

// GetRestPoints implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetRestPoints() map[string]sci.SmartContractRestHandler {
	return msc.RestHandlers
}

// Execute implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) Execute(t *transaction.Transaction,
	funcName string, _ []byte, balances chainstate.StateContextI) (string, error) {

	switch funcName {
	// consumer
	case registerConsumer:
		return msc.registerConsumer(t, balances)
	}

	return "", common.NewError("invalid_function_name", "function with provided name is not supported")
}

// GetHandlerStats implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetHandlerStats(ctx context.Context, params url.Values) (interface{}, error) {
	return msc.SmartContract.HandlerStats(ctx, params)
}

// GetExecutionStats implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetExecutionStats() map[string]interface{} {
	return msc.SmartContractExecutionStats
}
