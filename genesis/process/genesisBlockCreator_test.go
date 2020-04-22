package process

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/genesis"
	"github.com/ElrondNetwork/elrond-go/genesis/mock"
	"github.com/ElrondNetwork/elrond-go/process/economics"
	"github.com/ElrondNetwork/elrond-go/storage"
	"github.com/stretchr/testify/assert"
)

//TODO improve code coverage of this package
func createMockArgument() ArgsGenesisBlockCreator {
	arg := ArgsGenesisBlockCreator{
		GenesisTime:              0,
		StartEpochNum:            0,
		PubkeyConv:               mock.NewPubkeyConverterMock(32),
		InitialNodesSetup:        &mock.InitialNodesSetupHandlerStub{},
		Blkc:                     &mock.BlockChainStub{},
		Marshalizer:              &mock.MarshalizerMock{},
		Hasher:                   &mock.HasherMock{},
		Uint64ByteSliceConverter: &mock.Uint64ByteSliceConverterMock{},
		DataPool:                 mock.NewPoolsHolderMock(),
		ValidatorStatsRootHash:   make([]byte, 0),
	}

	arg.ShardCoordinator = &mock.ShardCoordinatorMock{
		NumOfShards: 2,
		SelfShardId: 0,
	}

	arg.Accounts = &mock.AccountsStub{
		RootHashCalled: func() ([]byte, error) {
			return make([]byte, 0), nil
		},
		CommitCalled: func() ([]byte, error) {
			return make([]byte, 0), nil
		},
		SaveAccountCalled: func(account state.AccountHandler) error {
			return nil
		},
		LoadAccountCalled: func(container state.AddressContainer) (state.AccountHandler, error) {
			return state.NewEmptyUserAccount(), nil
		},
	}

	arg.GasMap, _ = core.LoadGasScheduleConfig("testdata/gasSchedule.toml")

	ted := &economics.TestEconomicsData{
		EconomicsData: &economics.EconomicsData{},
	}
	ted.SetGenesisNodePrice(big.NewInt(100))
	ted.SetMinStep(big.NewInt(1))
	ted.SetTotalSupply(big.NewInt(10000))
	ted.SetUnJailPrice(big.NewInt(1))
	arg.Economics = ted.EconomicsData

	arg.Store = &mock.ChainStorerMock{
		GetStorerCalled: func(unitType dataRetriever.UnitType) storage.Storer {
			return mock.NewStorerMock()
		},
	}

	arg.GenesisParser, _ = genesis.NewGenesis(
		"testdata/genesis.json",
		arg.Economics.TotalSupply(),
		arg.PubkeyConv,
	)

	return arg
}

func TestGenesisBlockCreator_CreateGenesisBlocksShouldWork(t *testing.T) {
	t.Parallel()

	arg := createMockArgument()
	gbc, _ := NewGenesisBlockCreator(arg)

	blocks, err := gbc.CreateGenesisBlocks()

	assert.Nil(t, err)
	assert.Equal(t, 3, len(blocks))
}