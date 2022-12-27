package subclient

import (
	"encoding/binary"
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
	//[160 200 26 192 153 155 152 189 138 174 246 230 87 106 59 191 71 227 63 0 99 55 250 87 16 55 183 183 168 168 2 4]
	//14dp76EwTctDZmX8bgJV3jC6KsnCCpjwzvjMpm4tc2AkJN2L  0
	//efSnQr1KK8udzkovesB8TbRJo4QMXiyLbWEdiVzJH1RTXRZgE 1110
	networkId := make([]byte, 2)
	binary.LittleEndian.PutUint16(networkId, 1110)
	public, err := AddrToPublic("efSnQr1KK8udzkovesB8TbRJo4QMXiyLbWEdiVzJH1RTXRZgE", networkId)
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
