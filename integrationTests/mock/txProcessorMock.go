package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
)

// TxProcessorMock -
type TxProcessorMock struct {
	ProcessTransactionCalled         func(transaction *transaction.Transaction) error
	SetBalancesToTrieCalled          func(accBalance map[string]*big.Int) (rootHash []byte, err error)
	ProcessSmartContractResultCalled func(scr *smartContractResult.SmartContractResult) error
}

// ProcessTransaction -
func (etm *TxProcessorMock) ProcessTransaction(transaction *transaction.Transaction) error {
	return etm.ProcessTransactionCalled(transaction)
}

// SetBalancesToTrie -
func (etm *TxProcessorMock) SetBalancesToTrie(accBalance map[string]*big.Int) (rootHash []byte, err error) {
	return etm.SetBalancesToTrieCalled(accBalance)
}

// ProcessSmartContractResult -
func (etm *TxProcessorMock) ProcessSmartContractResult(scr *smartContractResult.SmartContractResult) error {
	return etm.ProcessSmartContractResultCalled(scr)
}

// IsInterfaceNil returns true if there is no value under the interface
func (etm *TxProcessorMock) IsInterfaceNil() bool {
	return etm == nil
}
