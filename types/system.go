package types

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type Extrinsic struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"amount"`

	Hash  string `json:"hash"`
	Nonce uint64 `json:"nonce"`
	Tip   string `json:"tip"`
	Era   string `json:"era"`

	Index  int        `json:"index"`
	Height uint64     `json:"height"`
	Module ModuleName `json:"module"`
	Call   CallId     `json:"call"`
	Event  EventID    `json:"event"`

	Value *big.Int `json:"value"`
	Addr  string   `json:"addr"` // other operation
}

func DefaultExtrinsic() *Extrinsic {
	return &Extrinsic{
		Amount: big.NewInt(0),
		Value:  big.NewInt(0),
	}
}

type SystemEvent struct {
	Phase  Phase         `json:"phase"`
	Event  Event         `json:"event"`
	Topics []interface{} `json:"topics"`
}

func (se SystemEvent) Parse(networkId []byte) (*Extrinsic, error) {
	extrinsic := DefaultExtrinsic()
	moduleName := se.Event.GetModuleName()
	eventName, err := se.Event.GetEventName()
	if err != nil {
		return nil, fmt.Errorf("get event name error:%v", err)
	}
	extrinsic.Module = moduleName
	extrinsic.Event = eventName
	if !se.Phase.IsApplyExtrinsic() {
		return extrinsic, fmt.Errorf("applyExtrinsic not success")
	}
	index, err := se.Phase.GetExtIndex()
	if err != nil {
		return nil, fmt.Errorf("get ext index error: %v", err)
	}
	extrinsic.Index = index
	if len(se.Event.Values) != 1 {
		return nil, fmt.Errorf("event values length not 1")
	}
	// todo
	// 目前我们只需要处理这几个模块
	if !(moduleName == System || moduleName == Staking || moduleName == Balances) {
		return extrinsic, nil
	}

	eventValue := se.Event.Values[0]
	switch eventValue.Values.(type) {
	case map[string]interface{}:
		value := DefaultEventValue()
		err := ObjToObj(eventValue.Values, &value)
		if err != nil {
			return nil, fmt.Errorf("obj parse error: %v", err)
		}
		innerValue, err := Parse(value, networkId)
		if err != nil {
			return nil, fmt.Errorf("event parse error: %v", err)
		}
		extrinsic.From = innerValue.From
		extrinsic.To = innerValue.To
		extrinsic.Addr = innerValue.Who
		extrinsic.Amount = innerValue.Amount
		extrinsic.Value = innerValue.Value
	case []interface{}:
		values := eventValue.Values.([]interface{})
		if len(values) != 2 {
			return nil, fmt.Errorf("length error")
		}
		addrBytes, err := ObjToBytes(values[0])
		if err != nil {
			return nil, fmt.Errorf("obj to bytes error:%v", err)
		}
		if len(addrBytes) != 1 {
			return nil, fmt.Errorf("addr bytes length error")
		}
		optAddr, err := PublicKeyToAddress(addrBytes[0], networkId)
		if err != nil {
			return nil, fmt.Errorf("bytes to addr error:%v", err)
		}
		extrinsic.Addr = optAddr
		value, ok := values[1].(json.Number)
		if !ok {
			return extrinsic, nil
		}
		bigInt, err := ObjToBigInt(value)
		if err != nil {
			return nil, fmt.Errorf("obj to big error:%v", err)
		}
		extrinsic.Amount = bigInt

	default:

	}
	return extrinsic, nil
}

type Phase struct {
	Name   string `json:"name"`
	Values []int  `json:"values"`
}

func (ph Phase) IsApplyExtrinsic() bool {
	return ph.Name == "ApplyExtrinsic"
}

func (ph Phase) GetExtIndex() (int, error) {
	if len(ph.Values) != 1 {
		return 0, fmt.Errorf("values extIndex length != 1")
	}
	return ph.Values[0], nil
}

type Event struct {
	Name   ModuleName `json:"name"` // module name
	Values []Values   `json:"values"`
}

func (e Event) GetEventName() (EventID, error) {
	if len(e.Values) != 1 {
		return "", fmt.Errorf("event values length !=1")
	}
	return e.Values[0].Name, nil
}
func (e Event) GetModuleName() ModuleName {
	return e.Name
}

type Values struct {
	Name   EventID     `json:"name"` // event name
	Values interface{} `json:"values"`
}

type DispatchInfo struct {
	Weight  interface{} `json:"weight"`
	Class   Class       `json:"class"`
	PaysFee PaysFee     `json:"pays_fee"`
}

type EventValue struct {
	//DispatchInfo DispatchInfo `json:"dispatch_info"`

	//other
	Who [][]byte `json:"who"`

	// transfer
	From   [][]byte `json:"from"`
	To     [][]byte `json:"to"`
	Amount *big.Int `json:"amount"`

	ActualFee *big.Int `json:"actual_fee"`
	Tip       int      `json:"tip"`
}

func DefaultEventValue() EventValue {
	return EventValue{
		Amount:    big.NewInt(0),
		ActualFee: big.NewInt(0),
	}
}

func Parse(ev EventValue, networkId []byte) (*InnerEventValue, error) {
	var from string
	var err error
	if len(ev.From) == 1 {
		from, err = PublicKeyToAddress(ev.From[0], networkId)
		if err != nil {
			return nil, fmt.Errorf("from publickey to addr error: %v", err)
		}
	}
	var to string
	if len(ev.To) == 1 {
		to, err = PublicKeyToAddress(ev.To[0], networkId)
		if err != nil {
			return nil, fmt.Errorf("to publickey to addr error: %v", err)
		}
	}
	var who string
	if len(ev.Who) == 1 {
		who, err = PublicKeyToAddress(ev.Who[0], networkId)
		if err != nil {
			return nil, fmt.Errorf("who publickey to addr error: %v", err)
		}
	}
	eventValue := &InnerEventValue{
		Who:       who,
		From:      from,
		To:        to,
		Amount:    ev.Amount,
		ActualFee: ev.ActualFee,
		Tip:       ev.Tip,
	}
	return eventValue, nil
}

type InnerEventValue struct {
	DispatchInfo DispatchInfo `json:"dispatch_info"`
	Who          string       `json:"who"`

	From      string   `json:"from"`
	To        string   `json:"to"`
	Amount    *big.Int `json:"amount"`
	Value     *big.Int `json:"value"`
	ActualFee *big.Int `json:"actual_fee"`
	Tip       int      `json:"tip"`
}

type Class struct {
	Name   string        `json:"name"`
	Values []interface{} `json:"values"`
}

type PaysFee struct {
	Name   string        `json:"name"`
	Values []interface{} `json:"values"`
}

type AccountInfo struct {
	Nonce       uint32      `json:"nonce"`
	Consumers   uint32      `json:"consumers"`
	Providers   uint32      `json:"providers"`
	Sufficients uint32      `json:"sufficients"`
	Data        AccountData `json:"data"`
}

func DefaultAccountInfo() *AccountInfo {
	info := AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Providers:   0,
		Sufficients: 0,
		Data: AccountData{
			Free:       big.NewInt(0),
			Reserved:   big.NewInt(0),
			MiscFrozen: big.NewInt(0),
			FeeFrozen:  big.NewInt(0),
		},
	}
	return &info
}

type AccountData struct {
	Free       *big.Int `json:"free"`
	Reserved   *big.Int `json:"reserved"`
	MiscFrozen *big.Int `json:"misc_frozen"`
	FeeFrozen  *big.Int `json:"fee_frozen"`
}
