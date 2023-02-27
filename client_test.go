package subclient

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"subclient/types"
	"testing"
)

const PLACK = 1000000000

var client *Client
var err error
var networkId = 42
var networkIdBytes = []byte{42}

var wsEndpoint = "wss://mainnet.ternoa.network"
var httpEndpoint = "http://127.0.0.1:9933"

//var wsEndpoint = "wss://rpc.polkadot.io"
//var httpEndpoint = "https://rpc.polkadot.io"

func init() {
	option := types.ClientOption{
		HttpEndpoint: httpEndpoint,
		WsEndpoint:   wsEndpoint,
		NetworkId:    networkId,
		NetworkBytes: networkIdBytes,
		WsSwitch:     true,
		Debug:        false,
	}
	client, err = NewClient(option)
	if err != nil {
		panic(err)
	}

}

func TestClient_scanBlock(t *testing.T) {
	height, err := client.GetFinalHeight()
	if err != nil {
		panic(err)
	}
	for i := height; i < 10000000000000; i++ {
		extrinsics, err := client.Block(uint64(i))
		if err != nil {
			panic(err)
		}
		fmt.Printf("len: %v %v %v  \n", extrinsics[0].Hash, i, len(extrinsics))
		for _, item := range extrinsics {
			if item.Module == types.Balances && item.Event == types.Transfer {
				fmt.Printf("%v \n", item)
			}
		}
	}

}

func TestClient_PalletInfo(t *testing.T) {
	metaData, err := client.GetMetaData("")
	if err != nil {
		panic(err)
	}
	info, err := client.PalletInfo(types.Staking, types.ReBond, metaData)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)

}

func TestClient_Block(t *testing.T) {
	block, err := client.Block(407)
	if err != nil {
		panic(err)
	}
	for _, item := range block {
		fmt.Println(item)
	}
}

func TestClient_finalHead(t *testing.T) {
	head, err := client.getFinalHead()
	if err != nil {
		panic(err)
	}
	block, err := client.getBlock(head)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Block.Header.GetHeight())
}

func TestClient_getHead(t *testing.T) {
	head, err := client.chainGetHead()
	if err != nil {
		panic(err)
	}
	fmt.Println(head)
}

func TestClient_OnChainApr(t *testing.T) {

	activeEraInfo, err := client.ActiveEraInfo()
	if err != nil {
		panic(err)
	}

	reward, err := client.ValidatorReward(uint32(activeEraInfo.Index - 1))
	if err != nil {
		panic(err)
	}

	eraTotalStake, err := client.EraTotalStake(uint32(activeEraInfo.Index - 1))
	if err != nil {
		panic(err)
	}
	apr := RatMul(RatDiv(NewRateFromBigInt(reward), NewRateFromBigInt(eraTotalStake)), NewRatFromInt(365))

	fmt.Printf("eraIndex: %v,reward: %v,totalStake: %v, apr: %v \n", activeEraInfo.Index-1, reward.String(), eraTotalStake.String(), apr.FloatString(18))

}

func TestClient_ChainAPY(t *testing.T) {
	// reward 123,004,589,200,024,601,923,222
	// total:316,727,875,901,025,829,290,136,572
	//rewardRat, ok := NewRatFromStr("123004589200024601923222")
	//if !ok {
	//	panic("")
	//}
	//toalRat, ok := NewRatFromStr("316727875901025829290136572")
	//if !ok {
	//	panic("")
	//}
	//apr := RatMul(RatDiv(rewardRat, toalRat), NewRatFromInt(365))
	//
	//fmt.Println(apr.FloatString(18))

}

func TestClient_Staking(t *testing.T) {

	//activeEra, err := client.ActiveEraInfo()
	//if err != nil {
	//	panic(err)
	//}

	delegatorMap := make(map[string]string) // validator  stash - > controller
	ownValidators := make(map[string]interface{})
	delegatorMap["5GCSraS9frHy8RQRKyMGbxCmRrVphVRFE2iQqDfhCkbd8taN"] = "5HdKwTY4YGEYSrPn9dt1aLRTqT6eexhBhJbGaxZh1vQqViRX"
	ownValidators["5FeBC65Gei67Lef3EoBTqfh1LFyZsFxAitRS3BUoyMz8rGeH"] = struct{}{}
	ownValidators["5GnSnCBzTyoLVSoQP1aXQLEmUBVyZt25z6J3e1MmUGzFyDvJ"] = struct{}{}

	controllerAddr := "5HdKwTY4YGEYSrPn9dt1aLRTqT6eexhBhJbGaxZh1vQqViRX"

	//accountInfo, err := client.SystemAccount(controllerAddr)
	//if err!=nil{
	//	panic(err)
	//}

	ledger, err := client.StakingLedger(controllerAddr)
	if err != nil {
		panic(err)
	}
	totalBondBig := ledger.Active

	totalUnbondBig := big.NewInt(0)
	for _, v := range ledger.Unlocking {
		for _, unLock := range v {
			unbondingBig := unLock.Value
			totalUnbondBig = big.NewInt(0).Add(totalUnbondBig, unbondingBig)
		}
	}

	fmt.Printf("controller: %v, bond: %v ,unBond: %v \n", controllerAddr, Big2Str(totalBondBig, 18), Big2Str(totalUnbondBig, 18))

	for eraIndex := int64(286); eraIndex < 287; eraIndex++ {
		currentEraIndex := uint32(eraIndex)
		eraReward, err := client.ValidatorReward(currentEraIndex)
		if err != nil {
			panic(err)
		}
		erasRewardPoints, err := client.ValidatorPoints(currentEraIndex)
		if err != nil {
			panic(err)
		}

		validatorPoints, err := erasRewardPoints.Parse([]byte{42})
		if err != nil {
			panic(err)
		}

		fmt.Printf("eraIndex: %v,eraReward: %v ,validators: %v \n", eraIndex, eraReward, len(validatorPoints))

		totalReward := NewRatFromInt(0)
		totalBond := NewRatFromInt(0)

		eraRewardRat := NewRatFromInt(0).SetInt(eraReward)

		for _, validator := range validatorPoints {

			validatorAddress := validator.Address
			validatorEarnRat := RatMul(eraRewardRat, RatDiv(NewRatFromInt(validator.Point), NewRatFromInt(erasRewardPoints.Total)))

			stakingCliped, err := client.StakingCliped(currentEraIndex, validatorAddress)
			if err != nil {
				panic(err)
			}

			validatorTotalBondRat, ok := NewRatFromBigInt(stakingCliped.Total)
			if !ok {
				panic("")
			}

			validatorPrefs, err := client.ValidatorPrefs(currentEraIndex, validatorAddress)
			if err != nil {
				panic(err)
			}
			validatorCommissionRat := RatDiv(NewRatFromUint(validatorPrefs.Commission), NewRatFromInt(PLACK))
			//fmt.Printf("validator: %v,reward: %v,commission: %v \n", validatorAddress, validatorEarnRat.FloatString(0), validatorCommissionRat.FloatString(5))

			if _, exists := ownValidators[validatorAddress]; exists {

				ownBondRat, ok := NewRatFromBigInt(stakingCliped.Own)
				if !ok {
					panic("parse rat error")
				}
				// https://github.com/paritytech/substrate/blob/bfc6fb4a95adccf96dadbe5e5e6bb8605d5a2c01/frame/staking/src/pallet/impls.rs#L86
				// 先从总的奖励中计算佣金  commissionReward = commission * validatorEarnRat
				commissionReward := RatMul(validatorCommissionRat, validatorEarnRat)
				//  ownReward = (ownBond / validatorTotalBond ) * ( validatorEarnRat - commissionReward )

				ownReward := RatMul(RatDiv(ownBondRat, validatorTotalBondRat), RatSub(validatorEarnRat, commissionReward))

				rewardRat := RatAdd(ownReward, commissionReward)
				totalBond = RatAdd(totalBond, ownBondRat)
				totalReward = RatAdd(totalReward, rewardRat)
				fmt.Printf("own vaidator: eraIndex: %v, validator: %v ,bond: %v, reward: %v \n", currentEraIndex, validatorAddress, RatToStr(ownBondRat, 18), RatToStr(rewardRat, 18))

			}

			for _, nominator := range stakingCliped.Others {

				ownBondRat, ok := NewRatFromBigInt(stakingCliped.Own)
				if !ok {
					panic("parse big error")
				}

				address, err := nominator.Address(networkIdBytes)
				if err != nil {
					panic(err)
				}
				if _, exists := delegatorMap[address]; exists {

					bondRat, ok := NewRatFromBigInt(nominator.Value)
					if !ok {
						panic("parse rat error")
					}
					// reward:= validatorEarnRat * (1-commissionRat) * bond / totalBond
					bondPercentRat := RatDiv(bondRat, validatorTotalBondRat)

					rewardRat := RatMul(RatMul(RatSub(NewRatFromInt(1), validatorCommissionRat), validatorEarnRat), bondPercentRat)
					totalReward = RatAdd(totalReward, rewardRat)
					totalBond = RatAdd(totalBond, ownBondRat)
					fmt.Printf("eraIndex: %v, validator: %v,nominator: %v, reward: %v, validatorBond: %v,ownBond: %v \n", currentEraIndex, validatorAddress, address, RatToStr(rewardRat, 18), RatToStr(validatorTotalBondRat, 18), RatToStr(bondRat, 18))
				}

			}

		}
		fmt.Printf("eraIndex: %v reward: %v \n", eraIndex, RatToStr(totalReward, 18))

	}

}

func TestDemo(t *testing.T) {

}

func TestClient_GetRunTimeVersion2(t *testing.T) {
	version, err := client.GetRunTimeVersion("")
	if err != nil {
		panic(err)
	}
	fmt.Println(version)
}

func TestClient_GetLatestHeight(t *testing.T) {
	height, err := client.GetFinalHeight()
	if err != nil {
		panic(err)
	}
	fmt.Println(height)
}

func TestClient_GetFinalHeader(t *testing.T) {
	header, err := client.getFinalHead()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(header)
}

func TestClient_GetHeader(t *testing.T) {
	header, err := client.getHeader("0x55725f9beb92ab6872161b5840c797976fe6d8c6cf528499a074a34a8b31cd5a")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(header.Number)

}

func TestClient_GetBlockHash(t *testing.T) {
	hash, err := client.getBlockHash(881)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}

func TestClient_GetMetaData(t *testing.T) {
	metaData, err := client.GetMetaData("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(metaData)
}

func TestClient_QueryExtrinsic(t *testing.T) {
	extrinsic, err := client.QueryExtrinsic("0x743095bbfd7c53256988fa43c5c64fc12ffaa80f2e2ef1ebd2b7ee3d5cdf7f25", 1712)
	if err != nil {
		panic(err)
	}
	fmt.Println(*extrinsic)
}

func TestClient_AuthorSubmitExtrinsic01(t *testing.T) {
	hash, err := client.authorSubmitExtrinsic("0x3d0284ff9c392efa851caa5a2152059f07a9fa978ff37fe9ef338eddf9db2955dd419d1601984afe2597ca2e68da26cbc5331e343a1dc719ff82d655d2b269b64aee25bf152f2bf5957b9dd25b32a23f3a15d275c9238bcbf4edf55a95a77cb162bd2a578a0004000600ff40fc14e721e231e8ff6106d2f07a0bac611714989c9f9de4c5f215b833e06beb070010a5d4e8")
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}

func TestClient_GetGenesisHash(t *testing.T) {
	hash, err := client.GetGenesisHash()
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}

func TestClient_SignedExtrinsic(t *testing.T) {
	//14fcLPyFxvPSv3mmGYmrNfg5Ln1otGbm8WeB7reeXwKPCb6K
	result, err := client.SignedExtrinsic(
		"0x186c09cac19834761b573b238b6542257d05b1fc5a57688311345d8cdf7e488d",
		"5FjKC4iC797yUWmFJuirEWqvVA2ABy3d41ugxZfHyrHs2AYx",
		"113000000000000000000",
		"4",
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	info, err := client.PaymentQueryInfo(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
	hash, err := client.authorSubmitExtrinsic(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)

}

func TestClient_GetStorage(t *testing.T) {
	storage, err := client.StateGetStorage("0x5f3e4907f716ac89b6347d15ececedca7e6ed2ee507c7b4441d59e4ded44b8a24213c2713e48b45264000000", "0x49f9877c851061f2076a2c0b8c7d37c93bcdb33a9a7ab8cfe20c8410d50437d8")
	if err != nil {
		panic(err)
	}
	fmt.Println(storage)
}

func TestClient_AuthorRotateKeys(t *testing.T) {
	keys, err := client.AuthorRotateKeys()
	if err != nil {
		panic(err)
	}
	fmt.Println(keys)
}

type User struct {
	Value *big.Int `json:"value"`
}

func TestClient_json(t *testing.T) {
	value, _ := big.NewInt(0).SetString("1234567891234567891234567891133333333", 10)
	user := User{
		Value: value,
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	u := User{}
	err = json.Unmarshal(bytes, &u)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
}

func TestClient_ActiveEraInfo(t *testing.T) {
	activeEraInfo, err := client.ActiveEraInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(activeEraInfo)
}

func TestClient_ValidatorReward(t *testing.T) {
	validatorReward, err := client.ValidatorReward(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(validatorReward)
}

func TestClient_ValidatorPoints(t *testing.T) {
	points, err := client.ValidatorPoints(0)
	if err != nil {
		panic(err)
	}
	addr, err := points.Parse([]byte{42})
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)
}

func TestClient_StakingCliped(t *testing.T) {
	exposure, err := client.StakingCliped(6, "1VsNodjV2hVhQNvLMtd2CSJ2naN6KF9Tat9JRfv3aRmquaM")
	if err != nil {
		panic(err)
	}
	fmt.Println(exposure)
	for k, v := range exposure.Others {
		fmt.Println(k, v)
	}
}

func TestClient_ValidatorPrefs(t *testing.T) {
	validatorPrefs, err := client.ValidatorPrefs(5, "1VsNodjV2hVhQNvLMtd2CSJ2naN6KF9Tat9JRfv3aRmquaM")
	if err != nil {
		panic(err)
	}
	fmt.Println(validatorPrefs)
}

func TestClient_StakingLedgerV1(t *testing.T) {
	stakingLedger, err := client.StakingLedger("esrSKbpVXVs4F3FaK9r2YJ5zyUVoQMj1PoieMriK9mwhncs9K")
	if err != nil {
		panic(err)
	}
	fmt.Println(stakingLedger)
	for _, v := range stakingLedger.Unlocking {
		for _, unlock := range v {
			fmt.Println(unlock)
		}
	}
}

func TestClient_DynamicDecodeStorage(t *testing.T) {
	metaData, err := client.GetMetaData("")
	if err != nil {
		panic(err)
	}
	raw := "0c980024080112201ce5f00ef6e89374afb625f1ae4c1546d31234e87e3c3f51a62b91dd6bfa57df98002408011220876a7b4984f98006dc8d666e28b60de307309835d775e7755cc770328cdacf2e98002408011220dacde7714d8551f674b8bb4b54239383c76a2b286fa436e93b2b7eb226bf4de7"
	palletName := types.ModuleName("NodeAuthorization")
	storageEntry := types.StorageKey("WellKnownNodes")
	decodeStorage, err := client.DynamicDecodeStorage(palletName, storageEntry, raw, metaData)
	if err != nil {
		panic(err)
	}
	var res interface{}
	err = types.Unmarshal([]byte(decodeStorage), &res)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestDemo001(t *testing.T) {
	res := `[{\"phase\":{\"name\":\"ApplyExtrinsic\",\"values\":[0]},\"event\":{\"name\":\"System\",\"values\":[{\"name\":\"ExtrinsicSuccess\",\"values\":{\"dispatch_info\":{\"weight\":132273000,\"class\":{\"name\":\"Mandatory\",\"values\":[]},\"pays_fee\":{\"name\":\"Yes\",\"values\":[]}}}}]},\"topics\":[]},{\"phase\":{\"name\":\"ApplyExtrinsic\",\"values\":[1]},\"event\":{\"name\":\"System\",\"values\":[{\"name\":\"ExtrinsicSuccess\",\"values\":{\"dispatch_info\":{\"weight\":2341737000,\"class\":{\"name\":\"Mandatory\",\"values\":[]},\"pays_fee\":{\"name\":\"Yes\",\"values\":[]}}}}]},\"topics\":[]}]`

	data := strings.ReplaceAll(res, `\`, "")
	var vales []interface{}
	err := json.Unmarshal([]byte(data), &vales)
	if err != nil {
		panic(err)
	}
	for _, item := range vales {
		fmt.Println(item)
	}

}

func TestClient_AccountInfo(t *testing.T) {
	accountInfo, err := client.SystemAccount("5HdKwTY4YGEYSrPn9dt1aLRTqT6eexhBhJbGaxZh1vQqViRX")
	if err != nil {
		panic(err)
	}
	fmt.Printf("MiscFrozen: %v, FeeFrozen: %v, Free: %v, Reserved: %v \n", Big2Str(accountInfo.Data.MiscFrozen, 18), Big2Str(accountInfo.Data.FeeFrozen, 18), Big2Str(accountInfo.Data.Free, 18), Big2Str(accountInfo.Data.Reserved, 18))
}

func TestClient_SystemEvents(t *testing.T) {
	systemEvents, err := client.SystemEvents("0xa92b8ade05727f6a30dfeee1ad1bff3f96acd104612cd30f2f646138e0bff46c")
	if err != nil {
		panic(err)
	}
	for _, item := range systemEvents {
		tx, err := item.Parse([]byte{42})
		if err != nil {

		}
		fmt.Println(tx)
	}
}

func TestWsClient(t *testing.T) {

}

func TestEvents(t *testing.T) {

	resList := []string{
		`[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":175517000},"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[212,53,147,199,21,253,211,28,97,20,26,189,4,169,159,214,130,44,133,88,133,76,205,227,154,86,132,231,165,109,162,125]],"amount":15100000280186545}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"NewAccount","values":{"account":[[160,200,26,192,153,155,152,189,138,174,246,230,87,106,59,191,71,227,63,0,99,55,250,87,16,55,183,183,168,168,2,4]]}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Endowed","values":{"account":[[160,200,26,192,153,155,152,189,138,174,246,230,87,106,59,191,71,227,63,0,99,55,250,87,16,55,183,183,168,168,2,4]],"free_balance":200000000000000000000000}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Transfer","values":{"from":[[212,53,147,199,21,253,211,28,97,20,26,189,4,169,159,214,130,44,133,88,133,76,205,227,154,86,132,231,165,109,162,125]],"to":[[160,200,26,192,153,155,152,189,138,174,246,230,87,106,59,191,71,227,63,0,99,55,250,87,16,55,183,183,168,168,2,4]],"amount":200000000000000000000000123456789123456789}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[212,53,147,199,21,253,211,28,97,20,26,189,4,169,159,214,130,44,133,88,133,76,205,227,154,86,132,231,165,109,162,125]],"actual_fee":15100000280186545,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":193890000},"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`,
		`[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":175517000},"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[22,4,115,170,186,199,218,216,54,125,240,140,253,154,246,215,216,121,137,126,102,16,124,59,227,55,205,131,178,220,74,97]],"amount":32800002627760505}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Staking","values":[{"name":"Bonded","values":[[[22,4,115,170,186,199,218,216,54,125,240,140,253,154,246,215,216,121,137,126,102,16,124,59,227,55,205,131,178,220,74,97]],199999850000000000000000]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Staking","values":[{"name":"ValidatorPrefsSet","values":[[[22,4,115,170,186,199,218,216,54,125,240,140,253,154,246,215,216,121,137,126,102,16,124,59,227,55,205,131,178,220,74,97]],{"commission":130000000,"blocked":false}]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"BatchCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[22,4,115,170,186,199,218,216,54,125,240,140,253,154,246,215,216,121,137,126,102,16,124,59,227,55,205,131,178,220,74,97]],"actual_fee":32800002627760505,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":2541723000},"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`,
		`[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":175517000},"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"amount":25900002110822597}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Staking","values":[{"name":"Bonded","values":[[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],199999850000000000000000]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"ItemCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Utility","values":[{"name":"BatchCompleted","values":[]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"actual_fee":25900002110822597,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":2024808000},"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`,
		`[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":175517000},"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[250,221,160,175,36,167,230,208,187,168,211,221,6,21,145,91,101,77,104,232,115,94,33,63,179,117,119,252,104,27,61,39]],"amount":11700001267154567}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"BagsList","values":[{"name":"ScoreUpdated","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"new_score":305548706344606861}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Staking","values":[{"name":"Unbonded","values":[[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],1234000000000000000000]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[250,221,160,175,36,167,230,208,187,168,211,221,6,21,145,91,101,77,104,232,115,94,33,63,179,117,119,252,104,27,61,39]],"actual_fee":11700001267154567,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":1181075000},"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`,
		`[{"phase":{"name":"ApplyExtrinsic","values":[0]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":175517000},"class":{"name":"Mandatory","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Withdraw","values":{"who":[[250,221,160,175,36,167,230,208,187,168,211,221,6,21,145,91,101,77,104,232,115,94,33,63,179,117,119,252,104,27,61,39]],"amount":11700001187943870}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Staking","values":[{"name":"Bonded","values":[[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],123000000000000000000]}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"BagsList","values":[{"name":"ScoreUpdated","values":{"who":[[192,222,123,127,102,140,201,28,244,159,162,14,135,2,142,92,196,54,23,229,231,146,245,80,190,36,145,25,182,27,102,49]],"new_score":305737785559574556}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"Balances","values":[{"name":"Deposit","values":{"who":[[250,221,160,175,36,167,230,208,187,168,211,221,6,21,145,91,101,77,104,232,115,94,33,63,179,117,119,252,104,27,61,39]],"amount":1983519}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"TransactionPayment","values":[{"name":"TransactionFeePaid","values":{"who":[[250,221,160,175,36,167,230,208,187,168,211,221,6,21,145,91,101,77,104,232,115,94,33,63,179,117,119,252,104,27,61,39]],"actual_fee":11700001185960351,"tip":0}}]},"topics":[]},{"phase":{"name":"ApplyExtrinsic","values":[1]},"event":{"name":"System","values":[{"name":"ExtrinsicSuccess","values":{"dispatch_info":{"weight":{"ref_time":1099929000},"class":{"name":"Normal","values":[]},"pays_fee":{"name":"Yes","values":[]}}}}]},"topics":[]}]`,
	}
	for _, res := range resList {
		var events []types.SystemEvent
		err := types.Unmarshal([]byte(res), &events)
		if err != nil {
			panic(err)
		}
		for _, event := range events {
			tx, err := event.Parse([]byte{42})
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.Module, tx.Event, tx.Value.String(), tx.Addr)

		}
	}
}
