package subclient

import (
	"fmt"
	"github.com/ChainSafe/gossamer/lib/crypto/sr25519"
	"subclient/types"
	"testing"
)

func TestNewStorageKey(t *testing.T) {
	storageKey, err := NewStorageKey(types.Staking, types.StakingErasRewardPoints, []byte{0}, types.NewTwox64ConcatValue(uint32(850)))
	if err != nil {
		panic(err)
	}
	fmt.Println(storageKey)
}

func TestPublicKeyToAddress(t *testing.T) {
	addBytes := []byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125}
	publicKey, err := sr25519.NewPublicKey(addBytes)
	if err != nil {
		panic(err)
	}
	address, err := PublicKeyToAddress(publicKey, []byte{0})
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
}

func TestAddrToPublic(t *testing.T) {
	public, err := AddrToPublic("15oF4uVJwmo4TdGW7VfQxNLavjCXviqxT9S1MgbjMNHr6Sp5", []byte{0})
	if err != nil {
		panic(err)
	}
	fmt.Println(public)
}

func TestHashHex(t *testing.T) {
	hex, err := HashHex("wohenidddd")
	if err != nil {
		panic(err)
	}
	fmt.Println(hex)
}
