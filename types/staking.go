package types

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type ActiveEraInfo struct {
	Index int64 `json:"index"`
	Start Start `json:"start"`
}

type Start struct {
	Name   string  `json:"name"`
	Values []int64 `json:"values"`
}

type ErasRewardPoints struct {
	Total      int64             `json:"total"`
	Individual [][][]interface{} `json:"individual"`
}

type ValidatorPoint struct {
	Address string `json:"address"`
	Point   int64  `json:"point"`
}

func (ep ErasRewardPoints) Parse(networkId []byte) ([]ValidatorPoint, error) {
	if len(ep.Individual) != 1 {
		return nil, fmt.Errorf("individual lenght error")
	}
	//if len(ep.Individual[0]) != 1 {
	//	return nil, fmt.Errorf("point lenght error")
	//}
	var validatorPoints []ValidatorPoint
	for _, ind := range ep.Individual {

		for _, individual := range ind {
			if len(individual) != 2 {
				return nil, fmt.Errorf("individual length not 2")
			}
			address := individual[0]
			point := individual[1]
			addressBytes, err := ObjToBytes(address)
			if err != nil {
				return nil, fmt.Errorf("obj to bytes error:%v", err)
			}
			if len(addressBytes) != 1 {
				return nil, fmt.Errorf("address addressBytes length not 1")
			}
			validatorAddr, err := PublicKeyToAddress(addressBytes[0], networkId)
			if err != nil {
				return nil, fmt.Errorf("pub to addr error:%v", err)
			}
			pointNumber, ok := point.(json.Number)
			if !ok {
				return nil, fmt.Errorf("point to str error")
			}
			validatorPoint, err := pointNumber.Int64()
			if err != nil {
				return nil, fmt.Errorf("number to int64 error:%v", err)
			}

			validatorPoints = append(validatorPoints, ValidatorPoint{Address: validatorAddr, Point: validatorPoint})
		}
	}
	return validatorPoints, nil
}

type Exposure struct {
	Total  *big.Int              `json:"total"`
	Own    *big.Int              `json:"own"`
	Others []*IndividualExposure `json:"others"`
}

func DefaultExposure() *Exposure {
	exposure := &Exposure{
		Total: big.NewInt(0),
		Own:   big.NewInt(0),
	}
	return exposure
}

type IndividualExposure struct {
	Who   [][]byte `json:"who"`
	Value *big.Int `json:"value"`
}

func (ind IndividualExposure) Address(networkId []byte) (string, error) {
	if len(ind.Who) != 1 {
		return "", fmt.Errorf("who length not 1")
	}
	address, err := PublicKeyToAddress(ind.Who[0], networkId)
	if err != nil {
		return "", fmt.Errorf("pub to addr error: %v", err)
	}
	return address, nil
}

type ValidatorPrefs struct {
	Commission uint64 `json:"commission"`
	Blocked    bool   `json:"blocked"`
}

type StakingLedger struct {
	Total     *big.Int        `json:"total"`
	Active    *big.Int        `json:"active"`
	Stash     [][]byte        `json:"stash"` // stash address
	Unlocking [][]UnlockChunk `json:"unlocking"`
}

func DefaultStakingLedger() *StakingLedger {
	return &StakingLedger{
		Total:  big.NewInt(0),
		Active: big.NewInt(0),
	}
}

func (sl StakingLedger) Address(networkId []byte) (string, error) {
	if len(sl.Stash) != 1 {
		return "", fmt.Errorf("stash byts lenght not 1")
	}
	address, err := PublicKeyToAddress(sl.Stash[0], networkId)
	if err != nil {
		return "", fmt.Errorf("pub to addr error:%v", err)
	}
	return address, nil
}

type UnlockChunk struct {
	Value *big.Int `json:"value"`
	Era   uint32   `json:"era"`
}
