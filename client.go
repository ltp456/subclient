package subclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"subclient/types"
	wsclient "subclient/ws"
	"sync"
	"time"
)

type Client struct {
	wsEndpoint   string
	httpEndpoint string
	imp          *http.Client
	ws           *wsclient.Ws
	networkId    []byte
	genesisHash  string
	debug        bool
	wsSwitch     bool
	timeout      time.Duration
	exit         chan string
}

func NewClient(option types.ClientOption) (*Client, error) {
	client := Client{
		wsEndpoint:   option.WsEndpoint,
		httpEndpoint: option.HttpEndpoint,
		networkId:    option.NetworkId,
		imp:          http.DefaultClient,
		debug:        false,
		wsSwitch:     true,
		timeout:      60 * time.Second,
		exit:         make(chan string, 1),
	}
	err := client.Init()
	if err != nil {
		return nil, err
	}
	return &client, nil

}

func (c *Client) Init() error {
	if c.wsSwitch {
		ws, err := wsclient.NewWs(c.wsEndpoint)
		if err != nil {
			return err
		}
		go ws.Run()
		go c.wsResp()
		c.ws = ws
		return nil
	}
	genesisHash, err := c.GetGenesisHash()
	if err != nil {
		return nil
	}
	c.genesisHash = strings.TrimPrefix(genesisHash, "0x")
	return nil
}

func (c *Client) Close() {
	if c.ws != nil {
		c.ws.Exit()
		close(c.exit)
	}
}

func (c *Client) ReadMsg() chan []byte {
	if c.ws != nil {
		return c.ws.ReadMessage()
	}
	return nil
}

func (c *Client) GetFinalHeight() (uint64, error) {
	hash, err := c.getFinalHead()
	if err != nil {
		return 0, err
	}
	header, err := c.getHeader(hash)
	if err != nil {
		return 0, err
	}
	height, err := header.GetHeight()
	if err != nil {
		return 0, err
	}
	return height, nil

}

func (c *Client) getFinalHead() (string, error) {
	var hash string
	return hash, c.call("chain_getFinalizedHead", &hash, nil)
}

func (c *Client) chainGetHead() (string, error) {
	var hash string
	return hash, c.call("chain_getHead", &hash, nil)
}

func (c *Client) getHeader(hash string) (types.Header, error) {
	var header types.Header
	return header, c.call("chain_getHeader", &header, []string{hash})
}

func (c *Client) getBlockHash(height uint64) (string, error) {
	var hash string
	return hash, c.call("chain_getBlockHash", &hash, []uint64{height})
}

func (c *Client) getBlock(hash string) (types.Block, error) {
	var param []string
	if hash == "" {
		param = nil
	} else {
		param = append(param, hash)
	}
	var block types.Block
	return block, c.call("chain_getBlock", &block, param)
}

func (c *Client) GetRunTimeVersion(blockHash string) (*types.RuntimeVersion, error) {
	var params []string
	if blockHash == "" {
		params = nil
	} else {
		params = append(params, blockHash)
	}
	var version types.RuntimeVersion
	return &version, c.call("state_getRuntimeVersion", &version, params)
}

func (c *Client) QueryExtrinsic(hash string, height uint64) (*types.Extrinsic, error) {
	extrinsics, err := c.Block(height)
	if err != nil {
		return nil, err
	}
	for _, ext := range extrinsics {
		if ext.Hash == hash {
			return ext, nil
		}
	}
	return nil, fmt.Errorf("no find extrinsic ")
}

func (c *Client) GetMetaData(hash string) (string, error) {
	var params []string
	if len(hash) != 0 {
		params = append(params, hash)
	}
	var metaData string
	err := c.call("state_getMetadata", &metaData, params)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(metaData, "0x"), nil
}

func (c *Client) PaymentQueryInfo(data string) (types.PaymentInfo, error) {
	var result types.PaymentInfo
	return result, c.call("payment_queryInfo", &result, []string{TrimSlash(data)})
}

func (c *Client) StateGetStorage(key, hash string) (string, error) {
	var result string
	param := []string{key}
	if len(hash) != 0 {
		param = []string{key, hash}
	}
	err := c.call("state_getStorage", &result, param)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(result, "0x"), nil
}

func (c *Client) authorSubmitExtrinsic(data string) (string, error) {
	var hash string
	return hash, c.call("author_submitExtrinsic", &hash, []string{TrimSlash(data)})
}

func (c *Client) GetGenesisHash() (string, error) {
	var hash string
	return hash, c.call("chain_getBlockHash", &hash, []interface{}{0})
}

func (c *Client) AuthorRotateKeys() (string, error) {
	var rotateKey string
	return rotateKey, c.call("author_rotateKeys", &rotateKey, []interface{}{})

}

func (c *Client) SignedExtrinsic(seed, address, amount, nonce string) (string, error) {
	runtimeVersion, err := c.GetRunTimeVersion("")
	if err != nil {
		return "", err
	}
	output, err := c.ffiSignedExtrinsic(c.genesisHash, seed, address, amount, nonce,
		fmt.Sprintf("%v", runtimeVersion.SpecVersion), fmt.Sprintf("%v", runtimeVersion.TransactionVersion))
	if err != nil {
		return "", err
	}
	data, err := checkFFIStatus(output)
	if err != nil {
		return "", err
	}
	return TrimQuotes(data), nil
}

func (c *Client) Transfer(seed, address, amount, nonce string) (string, error) {
	data, err := c.SignedExtrinsic(seed, address, amount, nonce)
	if err != nil {
		return "", err
	}
	hash, err := c.authorSubmitExtrinsic(data)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (c *Client) SubmitTx(data string) (string, error) {
	hash, err := c.authorSubmitExtrinsic(data)
	if err != nil {
		return "", err
	}
	return hash, nil

}

func (c *Client) RpcMethods() (*types.RpcMethods, error) {
	methods := &types.RpcMethods{}
	return methods, c.post(types.NewParams("rpc_methods", nil), methods)

}

func (c *Client) Block(height uint64) ([]*types.Extrinsic, error) {
	blockHash, err := c.getBlockHash(height)
	if err != nil {
		return nil, err
	}
	block, err := c.getBlock(blockHash)
	if err != nil {
		return nil, err
	}
	metadata, err := c.GetRealMetadata(height)
	if err != nil {
		return nil, err
	}
	events, err := c.SystemEventsWithMetadata(blockHash, metadata)
	if err != nil {
		return nil, err
	}
	extrinsic, err := c.ParseExtrinsic(height, block.Block.Extrinsics, events)
	if err != nil {
		return nil, err
	}
	return extrinsic, nil
}

func (c *Client) MetadataByHeight(height uint64) (string, error) {
	blockHash, err := c.getBlockHash(height)
	if err != nil {
		return "", err
	}
	metaData, err := c.GetMetaData(blockHash)
	if err != nil {
		return "", err
	}
	return metaData, nil
}

func (c *Client) GetRealMetadata(height uint64) (string, error) {
	//todo 在wasm runtime升级时，metadata在当前的高度没有合并，仅仅set code，是在下个高度生效的
	metadata, err := c.MetadataByHeight(height)
	if err != nil {
		return "", err
	}
	if height-1 < 0 {
		return metadata, nil
	}
	prevMetadata, err := c.MetadataByHeight(height - 1)
	if err != nil {
		return "", err
	}
	if metadata != prevMetadata {
		return prevMetadata, nil
	}
	return metadata, nil

}

func (c *Client) ParseExtrinsic(height uint64, extrinsics []string, events []types.SystemEvent) ([]*types.Extrinsic, error) {
	var tmpExtrinsic []*types.Extrinsic
	for extId, extData := range extrinsics {

		extCall, err := c.DecodeExtrinsic(extData)
		if err != nil {
			//todo
		}
		extCall.Parse()

		filterEvents, err := EventsById(extId, events)
		if err != nil {
			continue
		}
		hasSuccessEvent := ContainerSuccessEvent(extId, filterEvents)
		if !hasSuccessEvent {
			continue
		}
		hexData, err := hex.DecodeString(strings.TrimPrefix(extData, "0x"))
		hashBytes, err := Hash256(hexData)
		if err != nil {
			return nil, err
		}
		for index, event := range filterEvents {
			extrinsic, err := event.Parse(c.networkId)
			if err != nil {
				continue
			}
			extrinsic.Call = extCall.Call
			extrinsic.Height = height
			extrinsic.Index = index
			extrinsic.Hash = fmt.Sprintf("0x%x", hashBytes)
			tmpExtrinsic = append(tmpExtrinsic, extrinsic)
		}
	}
	return tmpExtrinsic, nil
}

func ContainerSuccessEvent(extId int, events []types.SystemEvent) bool {
	for _, event := range events {
		if !event.Phase.IsApplyExtrinsic() {
			continue
		}
		index, err := event.Phase.GetExtIndex()
		if err != nil {
			return false
		}
		name, err := event.Event.GetEventName()
		if err != nil {
			return false
		}
		if index == extId && event.Event.Name == types.System && name == types.ExtrinsicSuccess {
			return true
		}
	}
	return false
}

func EventsById(extId int, events []types.SystemEvent) ([]types.SystemEvent, error) {
	var tempEvents []types.SystemEvent
	for _, event := range events {
		if !event.Phase.IsApplyExtrinsic() {
			continue
		}
		index, err := event.Phase.GetExtIndex()
		if err != nil {
			return nil, err
		}
		if extId == index {
			tempEvents = append(tempEvents, event)
		}
	}
	return tempEvents, nil
}

// -----------decode ---------------

func (c *Client) SystemAccount(addr string) (*types.AccountInfo, error) {
	info := types.DefaultAccountInfo()
	err := c.GetStorageInfo(types.System, types.Account, "", info, types.NewB2128ConcatValue(addr))
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (c *Client) SystemEventsWithMetadata(blockHash, metadata string) ([]types.SystemEvent, error) {
	var events []types.SystemEvent
	err := c.GetStorageInfoWithMetadata(types.System, types.Events, blockHash, metadata, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *Client) SystemEvents(blockHash string) ([]types.SystemEvent, error) {
	var events []types.SystemEvent
	err := c.GetStorageInfo(types.System, types.Events, blockHash, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *Client) ActiveEraInfo() (*types.ActiveEraInfo, error) {
	activeEraInfo := &types.ActiveEraInfo{}
	err := c.GetStorageInfo(types.Staking, types.ActiveEra, "", activeEraInfo)
	if err != nil {
		return nil, err
	}
	return activeEraInfo, nil
}

func (c *Client) ValidatorReward(eraIndex uint32) (*big.Int, error) {
	value := big.NewInt(0)
	err := c.GetStorageInfo(types.Staking, types.ErasValidatorReward, "", value, types.NewTwox64ConcatValue(eraIndex))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *Client) ValidatorPoints(eraIndex uint32) (*types.ErasRewardPoints, error) {
	erasRewardPoints := &types.ErasRewardPoints{}
	err := c.GetStorageInfo(types.Staking, types.StakingErasRewardPoints, "", erasRewardPoints, types.NewTwox64ConcatValue(eraIndex))
	if err != nil {
		return nil, err
	}
	return erasRewardPoints, nil

}

func (c *Client) StakingCliped(eraIndex uint32, addr string) (*types.Exposure, error) {
	exposure := types.DefaultExposure()
	err := c.GetStorageInfo(types.Staking, types.ErasStakersClipped, "", exposure, types.NewTwox64ConcatValue(eraIndex), types.NewTwox64ConcatValue(addr))
	if err != nil {
		return nil, err
	}
	return exposure, nil
}

func (c *Client) ValidatorPrefs(eraIndex uint32, addr string) (*types.ValidatorPrefs, error) {
	prefs := &types.ValidatorPrefs{}
	err := c.GetStorageInfo(types.Staking, types.ErasValidatorPrefs, "", prefs, types.NewTwox64ConcatValue(eraIndex), types.NewTwox64ConcatValue(addr))
	if err != nil {
		return nil, err
	}
	return prefs, nil
}

func (c *Client) StakingLedger(addr string) (*types.StakingLedger, error) {
	stakingLedger := types.DefaultStakingLedger()
	err := c.GetStorageInfo(types.Staking, types.Ledger, "", stakingLedger, types.NewB2128ConcatValue(addr))
	if err != nil {
		return nil, err
	}
	return stakingLedger, nil
}

func (c *Client) GetStorageInfo(palletName types.ModuleName, storageEntry types.StorageKey, blockHash string, value interface{}, opts ...types.Option) error {
	typeOf := reflect.ValueOf(value)
	if typeOf.Kind() != reflect.Ptr {
		return fmt.Errorf("value is mutst pointer")
	}
	if blockHash == "" {
		finalHeadHash, err := c.chainGetHead()
		if err != nil {
			return err
		}
		blockHash = finalHeadHash
	}
	storageKey, err := NewStorageKey(palletName, storageEntry, c.networkId, opts...)
	if err != nil {
		return err
	}

	raw, err := c.StateGetStorage(storageKey, blockHash)
	if err != nil {
		return err
	}
	if raw == "" {
		return nil
	}
	metaData, err := c.GetMetaData(blockHash)
	if err != nil {
		return err
	}
	decodeData, err := c.DynamicDecodeStorage(palletName, storageEntry, raw, metaData)
	if err != nil {
		return err
	}
	//fmt.Printf("deoce data: %v \n", decodeData)
	return types.Unmarshal([]byte(decodeData), value)

}

func (c *Client) GetStorageInfoWithMetadata(palletName types.ModuleName, storageEntry types.StorageKey, blockHash, metaData string, value interface{}, opts ...types.Option) error {
	typeOf := reflect.ValueOf(value)
	if typeOf.Kind() != reflect.Ptr {
		return fmt.Errorf("value is mutst pointer")
	}
	if blockHash == "" {
		finalHeadHash, err := c.chainGetHead()
		if err != nil {
			return err
		}
		blockHash = finalHeadHash
	}
	storageKey, err := NewStorageKey(palletName, storageEntry, c.networkId, opts...)
	if err != nil {
		return err
	}

	raw, err := c.StateGetStorage(storageKey, blockHash)
	if err != nil {
		return err
	}
	if raw == "" {
		return nil
	}
	decodeData, err := c.DynamicDecodeStorage(palletName, storageEntry, raw, metaData)
	if err != nil {
		return err
	}
	//fmt.Printf("deoce data: %v \n", decodeData)
	return types.Unmarshal([]byte(decodeData), value)

}

func (c *Client) DynamicDecodeStorage(palletName types.ModuleName, storageEntry types.StorageKey, raw, metadata string) (string, error) {
	data, err := c.ffiDynamicDecodeStorage(string(palletName), string(storageEntry), raw, metadata)
	if err != nil {
		return "", err
	}
	res, err := checkFFIStatus(data)
	if err != nil {
		return "", err
	}
	return res, nil

}

func (c *Client) PalletInfo(palletName types.ModuleName, callName types.CallId, metadata string) (string, error) {
	data, err := c.ffiPalletInfo(string(palletName), string(callName), metadata)
	if err != nil {
		return "", err
	}
	res, err := checkFFIStatus(data)
	if err != nil {
		return "", err
	}
	return res, nil

}

func (c *Client) DecodeExtrinsic(raw string) (*types.ExtCall, error) {
	extCall := types.DefaultExtCall()
	newRaw := strings.TrimPrefix(raw, "0x")
	data, err := c.ffiDecodeExtrinsic(newRaw)
	if err != nil {
		return extCall, err
	}
	res, err := checkFFIStatus(data)
	if err != nil {
		return extCall, err
	}
	err = json.Unmarshal([]byte(res), extCall)
	if err != nil {
		return extCall, err
	}
	return extCall, nil

}

func checkFFIStatus(res string) (string, error) {
	status := types.FFIStatus{}
	err := json.Unmarshal([]byte(res), &status)
	if err != nil {
		return "", nil
	}
	if status.Status != types.FFISuccess {
		return "", fmt.Errorf("ffi fail %v ", status.Msg)
	}
	return status.Data, nil
}

func (c *Client) call(method string, value interface{}, param interface{}) error {
	params := types.NewParams(method, param)
	if c.wsSwitch {
		return c.wsReq(params, value)
	}
	return c.post(params, value)
}

func (c *Client) post(param types.Params, value interface{}) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}
	requestData, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if c.debug {
		fmt.Println("request: ", string(requestData))
	}

	req, err := http.NewRequest("POST", c.httpEndpoint, bytes.NewReader(requestData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.imp.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}

		return err
	}
	if resp == nil || resp.StatusCode < http.StatusOK || resp.StatusCode > 300 {
		return fmt.Errorf("response err")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	commonResp := &types.CommonResp{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.debug {
		fmt.Println("response: ", string(data))
	}

	err = json.Unmarshal(data, commonResp)
	if err != nil {
		return err
	}
	if commonResp.Error != nil {
		return fmt.Errorf("jsonrpc err %v %v", commonResp.Error.Message, commonResp.Error.Data)
	}
	return json.Unmarshal(commonResp.Result, value)

}

func (c *Client) get(path string, value interface{}) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.httpEndpoint, path), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.imp.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return err
	}
	if resp == nil || resp.StatusCode < http.StatusOK || resp.StatusCode > 300 {
		return fmt.Errorf("response err")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	commonResp := &types.CommonResp{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.debug {
		fmt.Println("response: ", string(data))
	}

	err = json.Unmarshal(data, commonResp)
	if err != nil {
		return err
	}
	if commonResp.Error != nil {
		return fmt.Errorf("jsonrpc err %v %v", commonResp.Error.Message, commonResp.Error.Data)
	}

	return json.Unmarshal(commonResp.Result, value)
}

// --------------------------------
//todo memory  待测试
//var cache sync.Map

var cache Cache = Cache{
	cache:     make(map[int64]chan []byte),
	startTime: time.Now(),
}

func (c *Client) wsReq(param types.Params, value interface{}) error {
	// todo reset memory
	if cache.LimitTime() {
		if cache.Len() > 0 {
			time.Sleep(c.timeout)
		}
		cache.ReSet()
	}
	id := param.GetId()
	sigChain := make(chan []byte, 1)
	cache.Store(id, sigChain)
	if c.debug {
		fmt.Printf("ws req: %v %v \n", id, param.String())
	}
	err := c.ws.WriteObj(param)
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	select {
	case <-ctx.Done():
		if _, ok := cache.Load(id); ok {
			cache.Delete(id)
		}
		return fmt.Errorf("timeout error")
	case msg := <-sigChain:
		if _, ok := cache.Load(id); ok {
			cache.Delete(id)
		}
		commonResp := &types.CommonResp{}
		err := json.Unmarshal(msg, commonResp)
		if err != nil {
			return err
		}
		if commonResp.Error != nil {
			return fmt.Errorf("%v", commonResp.Error.Message)
		}
		err = json.Unmarshal(commonResp.Result, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) wsResp() {
	msg := c.ws.ReadMessage()
	for {
		select {
		case data := <-msg:
			commonResp := &types.CommonResp{}
			err := json.Unmarshal(data, commonResp)
			if err != nil {
				continue
			}
			if c.debug {
				fmt.Printf("ws recv: %v %v \n", commonResp.ID, string(data))
			}
			key := commonResp.ID
			if value, ok := cache.Load(key); ok {
				value <- data
			}
		case <-c.exit:
			return
		}
	}
}

type Cache struct {
	lock      sync.RWMutex
	cache     map[int64]chan []byte
	startTime time.Time
}

func (c *Cache) LimitTime() bool {
	return c.startTime.Before(time.Now().Add(-time.Hour))
}

func (c *Cache) Len() int {
	return len(c.cache)
}

func (c *Cache) ReSet() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache = make(map[int64]chan []byte)
	c.startTime = time.Now()
}

func (c *Cache) Store(key int64, value chan []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[key] = value
}

func (c *Cache) Delete(key int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.cache[key]; ok {
		delete(c.cache, key)
	}
}

func (c *Cache) Load(key int64) (chan []byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.cache[key]; ok {
		return value, ok
	}
	return nil, false
}
