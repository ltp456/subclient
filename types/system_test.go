package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

type TestBigInt struct {
	Value *big.Int `json:"value"`
}

type Test struct {
	Value interface{}
}

func TestDemon(t *testing.T) {
	res := `{"value":1510000008629834011111111111111111111111111111111}`
	bigInt := TestBigInt{}
	err := Unmarshal([]byte(res), &bigInt)
	if err != nil {
		panic(err)
	}
	fmt.Println(bigInt.Value.String())
	data, err := Marshal(bigInt)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	test := Test{}
	err = Unmarshal([]byte(res), &test)
	if err != nil {
		panic(err)
	}
	fmt.Println(test.Value)

	data, err = json.Marshal(test)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

}
