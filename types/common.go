package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ExtCall struct {
	Module      ModuleName `json:"module"`
	Call        CallId     `json:"call"`
	PalletIndex int        `json:"pallet_index"`
	CallIndex   int        `json:"call_index"`
}

func (ec *ExtCall) Parse() {
	// todo
	if ec.PalletIndex == 13 {
		ec.Module = "Staking"
	}
	if ec.CallIndex == 19 {
		ec.Call = "rebond"
	}
}

func DefaultExtCall() *ExtCall {
	return &ExtCall{
		Module: ModuleName("UnKnown"),
		Call:   CallId("UnKnown"),
	}
}

type Address string

type EraIndex uint32

type Option struct {
	Value interface{}
	Type  HashType
}

func NewTwox64ConcatValue(value interface{}) Option {
	return Option{Value: value, Type: Twox64Concat}
}

func NewB2128ConcatValue(value interface{}) Option {
	return Option{Value: value, Type: Blake2_128Concat}
}

type ClientOption struct {
	HttpEndpoint string
	WsEndpoint   string
	NetworkId    []byte
}

func NewClientOption(wsUrl, httpUrl string, networkId []byte) ClientOption {
	return ClientOption{HttpEndpoint: httpUrl, WsEndpoint: wsUrl, NetworkId: networkId}
}

type FFIStatus struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type PaymentInfo struct {
	Class      string      `json:"class"`
	PartialFee string      `json:"partialFee"`
	Weight     interface{} `json:"weight"`
}

type RuntimeVersion struct {
	SpecVersion        int    `json:"specVersion"`
	AuthoringVersion   int    `json:"authoringVersion"`
	ImplName           string `json:"implName"`
	SpecName           string `json:"specName"`
	StateVersion       int    `json:"stateVersion"`
	TransactionVersion int    `json:"transactionVersion"`
}

type Header struct {
	Digest         Digest `json:"digest"`
	ParentHash     string `json:"parentHash"`
	Number         string `json:"number"`
	StateRoot      string `json:"stateRoot"`
	ExtrinsicsRoot string `json:"extrinsicsRoot"`
	Height         uint64 `json:"height"`
}

func (h *Header) GetHeight() (uint64, error) {
	if h.Number == "" {
		return 0, fmt.Errorf("get height err")
	}
	hexStr := strings.TrimLeft(h.Number, "0x")
	height, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint64(height), nil
}

type Digest struct {
	Logs []string `json:"logs"`
}

type Block struct {
	Block InnerBlock `json:"block"`
}

type InnerBlock struct {
	Extrinsics []string `json:"extrinsics"`
	Header     Header   `json:"header"`
}

type RpcMethods struct {
	Methods []string
}

type Params struct {
	ID      int64       `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func (p *Params) GetId() int64 {
	return p.ID
}

func (p *Params) String() string {
	return fmt.Sprintf("method: %v ,param: %v \n", p.Method, p.Params)
}

func NewParams(method string, param interface{}) Params {
	if param == nil {
		return Params{
			ID:      time.Now().UnixNano(),
			Method:  method,
			JsonRpc: "2.0",
		}
	} else {
		return Params{
			ID:      time.Now().UnixNano(),
			Method:  method,
			JsonRpc: "2.0",
			Params:  param,
		}
	}

}

type CommonResp struct {
	JsonRpc string          `json:"jsonrpc"`
	Error   *Error          `json:"error,omitempty"`
	Result  json.RawMessage `json:"result"`
	ID      int64           `json:"id"`
}
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
