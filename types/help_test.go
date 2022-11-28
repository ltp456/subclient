package types

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"math/big"
	"strings"
	"testing"
)

func TestAddress(t *testing.T) {
	addBytes := []byte{
		212, 53, 147,
		199,
		21,
		253,
		211,
		28,
		97,
		20,
		26,
		189,
		4,
		169,
		159,
		214,
		130,
		44,
		133,
		88,
		133,
		76,
		205,
		227,
		154,
		86,
		132,
		231,
		165,
		109,
		162,
		125,
	}
	address, err := PublicKeyToAddress(addBytes, []byte{0})
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
	public, err := AddrToPublic(address, []byte{0})
	if err != nil {
		panic(err)
	}
	fmt.Println(public)

}

func TestGabs(t *testing.T) {
	res := `[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":132273000,"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":2341737000,"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"amount":217184501}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"VoterList","values":[{"name":"ScoreUpdated","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"new_score":2001240000000000}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"Staking","values":[{"name":"Bonded","values":[[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],200000000000000000000000000000000000000000000000000]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"Balances","values":[{"name":"Deposit","values":{"who":[[109,111,100,108,112,121,47,116,114,115,114,121,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],"amount":173747600}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"Treasury","values":[{"name":"Deposit","values":{"value":173747600}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"Balances","values":[{"name":"Deposit","values":{"who":[[22,4,115,170,186,199,218,216,54,125,240,140,253,154,246,215,216,121,137,126,102,16,124,59,227,55,205,131,178,220,74,97]],"amount":43436901}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"actual_fee":217184501,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[2]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":818399000,"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`
	container, err := gabs.ParseJSON([]byte(res))
	if err != nil {
		panic(err)
	}
	children := container.S().Children()
	for _, child := range children {
		fmt.Println(child.Data())
	}

}

func TestDecode(t *testing.T) {
	data := `[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":175517000},"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[230,89,167,161,98,140,221,147,254,188,4,164,224,100,110,162,14,159,95,12,224,151,217,160,82,144,212,169,224,84,223,78]],"amount":15100000278168095}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Transfer","values":{"from":[[230,89,167,161,98,140,221,147,254,188,4,164,224,100,110,162,14,159,95,12,224,151,217,160,82,144,212,169,224,84,223,78]],"to":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"amount":123456789123456789000000000000000000}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[230,89,167,161,98,140,221,147,254,188,4,164,224,100,110,162,14,159,95,12,224,151,217,160,82,144,212,169,224,84,223,78]],"actual_fee":15100000278168095,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":193890000},"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	var events []SystemEvent
	err := decoder.Decode(&events)
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		tx, err := event.Parse([]byte{42})
		if err != nil {
		}
		fmt.Println(tx)
	}

}

func TestJson(t *testing.T) {
	type User struct {
		Value string `json:"value"`
		Name  string `json:"name"`
	}
	user := User{
		Value: "12345678912345678900000000000000000000",
		Name:  "test",
	}
	bytes, err := Marshal(user)
	if err != nil {
		panic(err)
	}

	bytes = []byte(`{"value":12,"name":"hello"}`)
	fmt.Println(string(bytes))
	type User1 struct {
		Value *big.Int `json:"value"`
		Name  string   `json:"name"`
	}

	user1 := User1{}
	err = Unmarshal(bytes, &user1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user1.Value)

}

func TestDemo(t *testing.T) {

}
