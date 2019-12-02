package sharding

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/state"
)

// Coordinator defines what a shard state coordinator should hold
type Coordinator interface {
	NumberOfShards() uint32
	ComputeId(address state.AddressContainer) uint32
	SelfId() uint32
	SameShard(firstAddress, secondAddress state.AddressContainer) bool
	CommunicationIdentifier(destShardID uint32) string
	IsInterfaceNil() bool
}

// Validator defines a node that can be allocated to a shard for participation in a consensus group as validator
// or block proposer
type Validator interface {
	Stake() *big.Int
	Rating() int32
	PubKey() []byte
	Address() []byte
}

// NodesCoordinator defines the behaviour of a struct able to do validator group selection
type NodesCoordinator interface {
	PublicKeysSelector
	SetNodesPerShards(eligible map[uint32][]Validator, waiting map[uint32][]Validator) error
	ComputeValidatorsGroup(randomness []byte, round uint64, shardId uint32) (validatorsGroup []Validator, err error)
	GetValidatorWithPublicKey(publicKey []byte) (validator Validator, shardId uint32, err error)
	IsInterfaceNil() bool
}

// PublicKeysSelector allows retrieval of eligible validators public keys
type PublicKeysSelector interface {
	GetValidatorsIndexes(publicKeys []string) []uint64
	GetAllValidatorsPublicKeys() map[uint32][][]byte
	GetSelectedPublicKeys(selection []byte, shardId uint32) (publicKeys []string, err error)
	GetValidatorsPublicKeys(randomness []byte, round uint64, shardId uint32) ([]string, error)
	GetValidatorsRewardsAddresses(randomness []byte, round uint64, shardId uint32) ([]string, error)
	GetOwnPublicKey() []byte
}

// ArgsUpdateNodes holds the parameters required by the shuffler to generate a new nodes configuration
type ArgsUpdateNodes struct {
	Eligible map[uint32][]Validator
	Waiting  map[uint32][]Validator
	NewNodes []Validator
	Leaving  []Validator
	Rand     []byte
	NbShards uint32
}

// NodesShuffler provides shuffling functionality for nodes
type NodesShuffler interface {
	UpdateParams(numNodesShard uint32, numNodesMeta uint32, hysteresis float32, adaptivity bool)
	UpdateNodeLists(args ArgsUpdateNodes) (map[uint32][]Validator, map[uint32][]Validator, []Validator)
}
