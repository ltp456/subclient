package subclient

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"subclient/types"
	"testing"
)

var client *Client
var err error

//var wsEndpoint = "ws://127.0.0.1:9944"
//var httpEndpoint = "http://127.0.0.1:9933"

var networkId = 0
var wsEndpoint = "wss://rpc.polkadot.io"
var httpEndpoint = "wss://rpc.polkadot.io"

func init() {
	option := types.ClientOption{
		HttpEndpoint: httpEndpoint,
		WsEndpoint:   wsEndpoint,
		NetworkId:    networkId,
		WsSwitch:     true,
	}
	client, err = NewClient(option)
	if err != nil {
		panic(err)
	}

}

func TestClient_scanBlock(t *testing.T) {
	//height, err := client.GetFinalHeight()
	//if err != nil {
	//	panic(err)
	//}
	for i := 13463845; i < 10000000000000; i++ {
		extrinsics, err := client.Block(uint64(i))
		if err != nil {
			panic(err)
		}
		fmt.Printf("len: %v %v %v  \n", extrinsics[0].Hash, i, len(extrinsics))
		for _, item := range extrinsics {
			if item.Module == types.Balances {
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

func TestClient_Staking(t *testing.T) {
	var networkIdBytes = []byte{0}
	activeEra, err := client.ActiveEraInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println("activeEra: ", activeEra)

	eraReward, err := client.ValidatorReward(0)
	if err != nil {
		panic(err)
	}
	fmt.Println("eraReward: ", eraReward)
	erasRewardPoints, err := client.ValidatorPoints(uint32(activeEra.Index))
	if err != nil {
		panic(err)
	}
	point, err := erasRewardPoints.Parse([]byte{42})
	fmt.Println("eraRewardPoints: ", point)

	fmt.Println("------------------------------")
	delegatorMap := make(map[string]string) // validator  stash - > controller
	delegatorMap["5CZaEUNfdFS2FsNQNiqct3c9BAaiQ1h1P69f98gZVVQFfFQf"] = "5FhWxkysbqck8EWce3FUuaMwUFnYWXBovRzsfU5Y3w9E7pMd"
	for k, v := range delegatorMap {
		erasStakersClipped, err := client.StakingCliped(8, k)
		if err != nil {
			panic(err)
		}
		fmt.Println("eraStakersClipped: ", erasStakersClipped)

		for _, ind := range erasStakersClipped.Others {
			address, err := ind.Address(networkIdBytes)
			if err != nil {
				panic(err)
			}
			fmt.Println(address, ind.Value.String())
		}

		ledger, err := client.StakingLedger(v)
		if err != nil {
			panic(err)
		}
		fmt.Println("ledger: ", ledger)
		address, err := ledger.Address(networkIdBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println(address)

		validatorPrefs, err := client.ValidatorPrefs(8, k)
		if err != nil {
			panic(err)
		}
		fmt.Println(validatorPrefs)

	}

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
	raw := "cc93ea00a5040000966d74f8027e07b43717b6876d97544fe0d71facef06acc8382749ae944ef02c0100000b93d72dcc12bd5577438c92a19c4778e12cfb8ada871a17694e5a2f86c374fc94000000b03b23766d70d0445943b290606521acaefee7660d521950faf2801c79d428c02f010000cec9b71b020bbfb41d2775d66d7911f22038ccae59c987108fe9a46d322a6e50fa0000021ba8ef466ccb7a06bfdefdba01817e620ada5954c34f617f6662e267dbda15a0fa0000040f4dbfe9b3333f7faf9304bf2ba9af0b4f850bab4e716759e0ebc8c4703c8c04fb00000507535074b7da9ae4989d5010cd0fcb3a33f8a516b721dd1a82c0a44901192e142c010006d2cf46cb11c18a8613ea1880c7cdf43ad1ac306a2a17aa4d5b4e4182b65c4214fa000006e11fd0d4df6c4765eb346aac47682cb7871da9ecfd235255f6eadb8392b20ddc8c000008a1e0292a6c1e6d0c4f0a2dc4d2c3f9b0e7c5d8213ce98d01b8e45be9edf15820cb00000aa041bf62d5b52adee07e7cbc0687ace86e77308a982c07db871ef2dab63901909700000c5e9de3da6f6dc71a0185f359418e0ef3177ceb8eb027d655807f00acbfa05594ed00000e518107b425628e49c4fee37fbc08020235d0664fc5052c4fa7415e5797c72c68c900000ec3f7291f82335606f98e16f5480b38b3da95d2e7ee9489a89a8a39f9dac956d8f900000edaa0d08e8b21d6e8f946eff381ad4be29aa63569a54a6a75c91b878b46303358ca000010b2f3fbaa08eb3361a1590bb8322b1fbef13bfd56bf5e04474d2a1ecaf9bb4894980000127a30e486492921e58f2564b36ab1ca21ff630672f0e76920edd601f8f2b89a309800001280a479ee3beca7af1636aca17582f30829782e2c9b1b9c72aaf8060563ab3740c90000128865fced1fe2ad870d0f1ef6ac3c73c78012dbaf73ee9db06ba403cb73a523d066000012ac5f0dfa73e4cf85d62a162b811cbc022a3d796aeb11cfe37fc1c555ef0a79e097000013be73ce92b712d00931d4980713bf4be8974255e1e51a7ed71ed2ea37f035ddd0290100161be097a4282aa58dae029624395fa7fc1c25cfa95aa7cc5100654963368442609a00001634f6ae2620175af504a54c459c3b40a7e2f8b89ea604d05c834a2e5d6d2b5778fa000016b94e2d5d12d60c7314cca383bf185ddda83f413da740a121601e3277d3083e6865000017316829c406a05cd9cdb8d5de5fb23d26b3672f8cbca1fcc6538833589a121afc350000182aed996417072b53ee48cd26d29fd7d986de458e1427e9fa9f21fcb261505b84ee0000189ef65d1a77acbc1492c95ea3a58cb70b9ffff211190368bbfa21198c734c2d089d000018a5657b88d667ab46d1070d58ecc1f15011f60df20f510ee0aac61754bb9e01082e010018f16f2628d93de560e2bdb6d1cb4aa3b681caad2c5ae723f26c09ecb676bc680497000018f6847e4d4c9bca7fdc694eeadfa8e36d2f49d581f3fe9e8599658f47dc2d5c0cbc00001a0e48b445798cbb689f86c9c17cd5428cec3bc370fe85a73cd3f1410c335d1e706700001a58ec699c897903d28984e42202dc216e2e8f7023c846926179c0c38562e4e0886800001c1975391b6d548727640335c328519ed843487c5e2a9e1a5fda09a94590bc51f0c800001c79a5ada2ff0d55aaa65dfeaf0cba667babf312f9bf100444279b34cd769e497cf600001e86387b4f6975e4cc574038cafddc0fb14d09b846da3c657e27599cf2aafc51349900002076db3b59f25dbe141f2f46fe0208978de4f9a3929376253b15968b55025e487cf1000020ac6c23e69518f5c048cdd4341f431d23f1bdcba3abfaf7349241db61ce1317289b000022abc926b7b1a5a0a202183919de1089d04939e09833ca7536b52cce9fa9bf3528c8000026a81cc7f1e72380949491cb9538d125d40de48e631e0e8bc40964fddee59bcc10d60000282a194090fd6715e06430d8a6e9c682f021eaf398830b10db94ca8c27c9ae4c30980000288197fe1d9f6b12ee7f9656663791366fd5ea4ec01999d7a7f8b84e22a0591bb0cc000029a843ce5a67c2d285cc5f302db173ee4563e7af383dfc6b9a8fa64e148efd12a42d01002a9c511c5d5ce66c809bc497eee0e5dacaadac372d56fb12ff8361b8f42c0760f0cd00002ae939fcc0d4fb0bf095e48ed294463613a1e30381d269940ceafcdfaf87ec94101c01002c2a55b59a282fb173b7faf6e82984686f6fb83b0293dd498db6329464a45d1cc89600002c2a55b59ed1954e80bf1349ee5882d429d261cf9962bc5d88a1fc176e60c91870fd00002c2a55b5a0ac2e467c72c205ae9252a4e8a0c206703950f095c31ea3f022493a882b01002c2a55b5a609baff13899d4ba4bafec105038d66a716494968fae1a849d2dd5ac43001002c2a55b5a6413894e13836fb0165e6adce7d77c06dccf42b3b288397a27ddd3b48cb00002c2a55b5b05723585c7421428d9ef451313e33e13803424b0f8cabd383d0df37609a00002c2a55b5b7e13a772e0b693c3b351d2fb5e5b4da18ac379ebdb2f1f2e7559776f8fc00002c2a55b5c60139f9572ec1f88f541744ad98ca9cd09b5641e84088ab0ab3920cd42f01002c2a55b5c69b5e131fb0f65ac7ca707f4bc53e4d991a2d1971ab5e702f69f45c4c6800002c2a55b5c7e135ae443c6ed951338d4b32be6cf61ac5bd000b7a11c26b307820b0fe00002c2a55b5cf8d00511ef54ac7a60773810e906befb1b322f2d58199963dc97307a02c01002c2a55b5e0d28cc772b47bb9b25981cbb69eca73f7c3388fb6464e7d24be470ed4cb00002c2a55b6068e3b8626608321044a89b82fc4898ece34524659f48aa72aef556c98f800002cc16da9d1f7271475075aa8eb5c6667714426b8c41dbecf92bdedfa462b7163b02b01002e2409b5ef509e1e584584edc945545f42fcbb3f288f3355e9194206b4ce773fe0fb00002eef2aee654d4975535f2701af86ba6d169c2c9a1599b16635a2a5e4640db94d80a7000030198db6b0a97949715f03d1cdc7fedd6d29b22cadc8d639ea0bcd2cf68458295c940000303349fd10535499bffaa9145a99b0e9e28df29c3ff9ef1a1baa3be64099c77c14c80000305c5105935cf62fb0c9408e919d4855b7c6c6972491f88261835cc4aaf602453499000030606b4c1b89b4e562efafe76bec80154ac8b3e16e04c2e0f619bcdc0a5eef5210cc000030cfdb48ff7f33b08499dfc618a8ef9699b8345fa65f0b1339eb8eec3c0e45559899000030e3f26094a2536e02b6879a2a92752493302d84029b80ff449eba8b9dd80d1de8c600003296c9b3a6546d2319a764e00e1126215b776cb27571db2ef0392bbfbc66d45f70c6000032f18f75e68ce71cfde8551a80a9aeaf56afe63c1773fafaeca7260f8132a2144c900000341646b4f31605eccaad023a8e17d6d8ac423420700db4bf0f125809c682e63f842a0100345d01c8b0b8d24ae3955161f24a3c3f832405539b54d997501eeaa68c52521d1c980000346ebc3380be6816f828d1d1df372c51fbe99c95a321d7510403bb98f067695ed8f9000034968759c88695ed60fefe374ff71252cedf3a49f238b5f9ea4442afa1184f2968c90000360d52090bbbd32b598beb80197c95c1004505135f0493b089334a52e3a1814098cb0000361b17b51e238769cf7ad4fd89ef1bb22c3c916c6c760ac57ea6c2d70886096fc4c70000364c29bfbc9f06a42b5cf37ffd831e91c843cc25d8b90071546810ecf279e45884cb00003694566bfebcfd72c7816a87f2ccf2963c0c794737055ead20b502a89432b9006c66000036fd9b64e99689363368298680cb35750a594103048bffd839af770fe5536c6cf0c80000380766b2d93f1737c6d8dc713e4a6d69b961a14a11f676fc479b1ebbe4b81f03a096000038a295559d8977464fd8cdd133f8805f2388e42a6e009219247048a27d9ac06bbcfc000038c95b31fe81ba229b2e2216f0bf6812ca45dc3ad5c49e15b48e0eef5555bf40bc98000039a1f89df89766e6e0771c1cf047fb712fdaff18eea966267c01ff74969876f6086600003a5a7e1a55916817a23b08433bff85bafc28855fbe67f8d2952b6b5b281640304cfe00003c017930b46ab5a4413bf3153b001287ed5ff7fdbd2734cf69abc843f4ee044794fc00003c1f7d137511560f34bd17ff8e9e48a9831f24f71512ce96c39267e8234b5e26f4ab00003e5dabd7d2d5758f94d780e23a088a9e793c6d8f8229efd810f2be2e966fe332a8fc00003eea2e1c17fc878e20fc5f88c82d6b08fd3324cbd24c176351aba2b06d874f21d8f90000408492998608540c567944fa19068c7a929675126a515c02fa15e8d50292813c94ed000042019132c3ccad3b962a76942099bf218649db3f0bd6b5336d94d5c8abdfb16280ca0000428f5be8df55fc2d1a12eccbf00ad2cb9e7598e342a58845ba37fb1749764e43a497000042efa2e57a813989da4bac4551e4010ee45003fc3f360f5202a958b2b1a299186865000044882dbb3ffba5fdce03b87b0e750f293fcf840964a1d5fec119e89fee2e0e42041e01004687fe7e263038bbc1b371c67e4f18b23a6b2f65c9cfe87675db0693eca8161d70fd0000484cdc76e0b6b2cb4e30850327cf37e717d91e343a62bbfaded38aa8133cfe3438c7000048931d40a4936659a5b03e419db8d3dfc8fffaa9737ad5771651552cc840210a64cd000048e181bb690336782cfd0b422faec0a561b5b35a6c2277f8073127ea8cf67b35402d01004a8de2c6b6c1690a27aaa4cb64e0168bced54ab568958965187646f06442be6d149b00004c260b3a7196660071ca1b27e79bf234cb7efaa125300a696a0a42ee686b173400cd00004cd38181e02e880a53114ccaf987a6f51a10bc1d172d509bab6f7e9d6eb2e00bccc90000509973649092168d075f8ff9192b46dc74a515dd2bcc03a1061763c396e5dc62242b010050c76f31fd1d42d7af02d5349ff5adc86a326aa8ffaf3f651f21dd792e061361a42d01005606fcd617bb7328f828ebc07396ab27f038eb5f036125a2ab5b60ab9218f24284cb000056196c14df0a7036b943ecd01396685e799f786c0f131796c06850ec9342ff0174cc000056ad64b58553e04382fbb95145efe3dfda88b5512c0f8fcf1cec8ee0a05dd172bcfc000056f462e2a9afdaae6600562a4d94c4eed98ebfbabef7785a111dc29d46f6e11894f7000058a96f0508bc3ff4fa8305c6e96a578dfd648a55aef61fc74291ea67390bb56054c9000058e0e1940ab2089ce4f65dc63abad1a6b654561a9a0b09fe5cc679f0ded3fc28889a00005ac7f6af5aeb5364188840d02f0e74e813e6d9cc0398d6994b66727658a4fb305ccb00005ae7010248daf19a0b83b3d131f63a693785222293af4354035b8dce851fb02bb0f900005c3c4c795e554a782e59ef5043ca9772f32dfb1ad7de832878d662194193955ee4ca00005e10fce6eaa355c47476dffce1cf118a1f0c28bc4209c8242d4271d07bfc2c03249a00005e1eb942e5f591bc1d902fc6a04dd7ecadf5f4916f2597df0505c6c521412c4bd8c700005e348817abb98cb962fc0780a47ebd471d9c318395fa80b4529a64cfabb2e32c30fc00005e4706e0f10c8d2a1cf27d967c5945b8f4e5d7d164212915a702415776e2947a0ccb00005e5d733b1d460c0836c2f01044537e6846e6200923c2a7e1e4cbea4dff44364744fc00005ec3359163109e793b726065a0ace5a201b72fbfc2d346c181773a04ee3ec87248f800005ed6f4b68a60a117a32059e96678f22fa086c505f7a8ece13c7a2e78b4788f1500cd0000600240b3b79fcd002d75c09a7f251b57429449c646bea7a72403a051fa8a403604c90000605fd1308af1ce85bab5ba3fb19b330ab7dac29e01ad501420560f44df7e0e1c9cc7000060a8308dfcb6611d34c0780562713dddb62fe542ee0d51f3f24aa5ad776b2e52e429010062ff35ef068f1845fc504efb7757fbbb7affa0b82e3fa8d2b6432dddf2e312239cfe00006558e62bb3e9e7da178a13f657c9a0e86044335a77a672f9f571489f43ae07b674cc000066207ff50019191490dc838fa46d1659070c9a4c3e739343ea009c38fb51cf6e34990000669cceb7c4d0fe56de471aad2b528b7229aec90b8f0096956de8d53b8aa87a0094c5000066dbc4ee5444f2e09c0ebb745c230a8dc516302a8d28a42cfd7882cdd8162e7928ff000066e38a65c3f53a195f63cd734c5839433d2679775f6630687bdc07043fd4a11570990000688140b645221a0239782a1f27f7dfa3b2967ba30b52bcd2c922d993cff3d27d8c91000068acd1796e80866ed99e01d0afb1686a365b3b228829e4d227eac2a9cd324d2ed0cf000068bcfb86663abec620784eec51572fa12ba5e10a7f312981e4d6f1bce379b939a46500006a0da15516a63ecc95cccffecec2a2aed962e92cab661af4b1f76b65e429cf7938f900006a13a77ddc9b38c10586c94c2928fa1af377757f877984ee76c511c6079b2f69e49300006aadafc32a25db2698407783d32e77416235046f9c12e05e2fc654682cf04c72589800006af358e5650b61943ba709efc5dbc501405e04d5de5798087d6c727027511a65c4fe00006cbabca2990dd4eb63081dab625a12cc535aa087ed623b0405becabfaf5b4d3200ff00007046b09e1b89a70d6801f8f3ea0fb3c19e8b7bbab64617b312be826fb08a590e589300007091f937fba948654220a41ede536b0d62cc30d20274a28005b8026564db8d2598990000729324ff6798093939a73546e0f3d53a9cd7d4e938d238145c9422ce9f0beb0750f50000741581af3c9bff5769bde735ee623f7580577f4dec888c9318485e16b6a0d81034990000744d9a778b5c53b54eb743dd93cc90079261a6e7fdffd8798e6406817d125d7b24c7000074c92d1501b154b7290f4ac36e45176d4c6bc2548cad5b254def72c3c2773e5a2cc900007638bf10d5a2077394a658561a0900bbb1359ef0510e411d5edbb168bd59b75314c800007677c2a4deb2689377319a8e830d5d9ce0ae32b95f096a7447f8f711b7f0333f08ca000076c26a1fb9acbdd56be00d4c44901856929b9d2a879caad6119ad0417e99494910cc0000783770544556ce47aae6e9e9c6cc32bb84a3df7ab11a5d30b0618e00e109d66c40fb0000785e826e5bf7f587cbeba82493e55bcc159d3d21119838e106abb7091ac4422df0c80000788abef41947f200a91aef46820c5f4bd89438da3ac0090bf25d70b90a3e411b48c600007a8b406e98e5b8439e8654cb022c4d61fc458de099d54e79699b28f1fab01b5b2cfb00007a8f8779b66af64a095e9e6c7600a89d7c33a77fb95d5b503505d264ad559531c4ef00007ac0c7e610ba022e2a048fbf5fd9949f96d813815f7174ce30f0124bc52c96644cc700007c2754f11fe0d5e6b3c2174eacd1fe36d2738bfbefb541fb4f56eb298dc1491a88a400007c3f190f0abecf2b39643a21a67df302a024487d84128d5bc68fcc445ac23a06e8cb00007c5371a38473955e88760d8f63c507b53e5dbda9205a524cc8d3b2c2af028f027c8d00007c6e71dab02a17292e9c2f3fb01d12a3895af33971a8fa710bb2820cb9d7a31ee06500007e8da7a268f01cb9d9a5307295a9752e292c4f715cf9e6b46d542f161cfe2823706700007eb71c2356edd039e714141ff023579e979466c0f3acec566714ebeed15f5d6c58ca00007efd4d272e412f7169470d93fffb03075854afeec6f0230ad4530f1d55fc48478895000080bb486f9b44cd9c3200aa5f85cc51e8921abc81bb79514248e4c486e85cf965b8c9000080e58bb10009a336451ccbbdc1c5aa99406bb0c90af50cb09daffef49b73c10cacc6000082635369d20e36ccce82c98283410bfa4a4153e31bedb52b33cbce5543cba21bf49700008440a356a35192f4653f72a11d37c09064e5a2ffb905403507570de01b1d1462609a0000844c2eea62566660b545e5cb75c4ffa2412b8a93bea1734156a8b10250ce88512cfb0000845d26853b5caf8444b5995c50b7b9bab9e7ffb27cd20641cc4c9f04e5942b4110cc00008639d9f577d4910627cd5d4643b205fc5476962afdd1d0068cbdd01f57618f6018c9000086627eb7d23d90c3a940e43ed8a705bb2845911ddcef53bf90560ee0b673e329d4c60000880678543460d4ad15fb10ea9ec36ab97d47e9317b7291fbe5c0c69aff6dab6430c50000886aa781773b34525e005ff2f4372d167aff28f24bf07863d087c088647dcf3bf8ca00008a32f59713f0a129fbc395dbc853f51ab53d45d1684c4bc8ddad89fd55fc096f106800008a449fac4875a24774175651ffcacc851939eff231fcd2e270e1e2942990205960cc00008a64dca67c64d9361901a1415bfb3469b000d0bf7f1d439824cec71f870221590c9900008a650412c92229bf1a2eb1c0cab9c133d54c5b82cc723b202ee634e925effa6afcc600008b1789281b38392a08ab0516a21dee49fba6f5e55381106ff7ae190563cc84f4e8fd00008b7602f67d9d61682964f3b5989f357752f75dfba430603da1384059b79f1535789600008c076c2dd6e9caa8d135ce6472493ea335eaae5e8250a71a2365929b99216025a89300008c23324b0cb29e4fd1a68cb08febe58b50e39d8afdb5f752d6c26c8ba52fc00268ce00008ccb33f77fbf009ad78e5e7d2b36bbf5355b9c4ceee22d75b1323c0c196e1c39c0fd00008d6bde6dcbc00b588a0ecb71f4c516d9ede311d3fcf942f89228221b92a2c1a328fa00008e03d1d653c37a5c2dbe0ff4b96d7bb07c9908594a5df84740b206cc51111776606800008eaf11411ed3bda42c77d9c9cc71f5e4a9bce979179446e9c8aa0809581e432c389a000090061acc3eb19871101b732f5fab60ff4e139dde0ed27b85893ed15ace3cdf1a30ca00009088b991261091d43c92ce790d4c8fb007b7e8d19c9c619a7d722c6fbdb55f04d8f900009208af2fb9f1511facbd517ba7ec0296d6fc9fd010896eec0f43a415b68b3706ccc90000925210734616021729765502d6f0ad63ec825f88cb725511c580ebbeeca45d7bc4f900009616812063edb4602413868c7b859655016342af7d6a64931e4520aef65a4627f0fa000096625a0cbd0931ad831add3bcaa6320950385aec23b3854c6ce987de1c9f8837fccb0000996a48be9e514c627122d28cbc8128389b88b230a85228a2fd091cc5d8292b39c4bd00009a71f753f8576834b7b43f9c33b61c177480e31c4b55308297fe3e7dfce06f7b349900009c665073980c9bdbd5620ef9a860b9f1efbeda8f10e13ef7431f6970d765a25780f200009c6a3401d06cef30fdbb33901328f3611dae8253708779a5d66179c967582635e89900009c6bae578812080b2f9ecdef913f4f68abc3dcee697f505625e8430c8ce3122c7cc900009e2e20b4c76b2d6b9210492cdbeb9e56c04aa9e81fdfad02f67e05c1fbe6401260b800009e553a94c461ade68e689e1372c3c863d1b49dbb629954cc62c263390c3a4a6e109000009effc7fb904ca3d8cc4f45b464ee4dac705c89888d298bc9dfd2ba563d5b3e3f80fc0000a0209c10c7d6633de64e166a68d9b45e7f02aacd9497f5c552e577abcdb17d2814ff0000a09b87b34880c5375ecca849557dc87a00a6243938d5882017fa0d1f60193815b4fa0000a28942a9d2e8c8501860b847eecedc45d602e614b8b0849b959607d0dec3d0710cfd0000a2c0af65e4f643cc81978d9960c6f25b4234e61bb04dd91a72a28eec71fc2437588e0000a2fafeae641e6e264d77723c00ab05f503db48ca3597cb3242c2b54d90abd01de0c90000a4a2b4ee5ad0bc423c3e2b5492545f2096550ff497b0749d535b885bd96d1828acf80000a4d30a367570be5fc2ed7786254b990be8e9e4b0de0ccfa6ba610e445cbfa826f0960000a5621bacecdb9d44ee96eb79109d23fb38f7291f9caff0d3dc1467c66847dae7a0be0000a685d6e1c48f7fe940f890577a57f1e6c131ae624b01d63867e2462012a0d945bcca0000a81b8a2a03afd92af18f0338e43afc504c6018aca6d9197c0d3a149ce65efa0eb8c90000a85ee90d7e89f7dbc08ea9b625e2efae557cbccb1c4337645cf5ea7566db661540c90000a864be32055e98c2b3216a463e4f3a8d345fd46f1120570d7f59276c464b77c3bcc00000a884d4694bef6d1c6901475028e23f24aeb0c769e7a885441294b387295e6b04ecf90000aab69386e2b8264fbdbe9ab35706cee70fe33147d42ebe2b45b0521b43831d54689c0000aab8591c28b705b3b8b67fba34373847994b66c060dc32b225139544ccaa4335b0c70000ac12e1c84e6b11175dcb8ee1f13f52d366ee23781f2042020dda2653efdbe84be4660000aca71966867cbad442367cb608c3088dc1176a1586d6efbb5137e579e2162248ccc90000acb208c5e0fed833c7a84dd41f8c7922179266fa7ce4dfa9604bbbb2d650fd1a68c40000add1d40ef104fcb24be637eaeb9704064663df235fea0799829e270333774f3bc49a0000ae0900005ec8554ccc01c89bc73d0d45a4f307fabeb6cfe87b141e246e33e756f8c50000aee72821ca00e62304e4f0d858122a65b87c8df4f0eae224ae064b951d39f610dc640000b023d129d9a0cb9490d097dbd3ca947d4830d3a6d7e0fa9975ff2789d9d97352c4cc0000b0a4b1e8778ee01886d8d3f11969d7ea564ce53c7f36dc59d112537048b0b01ee4ca0000b0c9266ebf0958d0d433f9d06864cf6e5a4f615ab978db196e6a9413ec2de4489ccc0000b0ea8d006c28df8d2833defa4850a6b5ff4717c785d481603d2abc442005c65a94fc0000b2c4d2772a169e3ccd912db6415a39a6f0228c89c7045556951b779133062375a0960000b65273df489c774e8c4da872b50219a8cdd40cd9e22432320aabac17c663d262a4c90000b6979b7897a0fbf4d194190973a6487277b671932c68712c8861a3c0945a861e54f60000b81712e18e76cfe727bda3fa65662405f584f1ed059aac110ab6f41851ba7e2624ea0000b8473b5e37e993092f313f54b58762953589245f5552612ec90042a6d59d807ae4ca0000b87cf5d33beb907cb0384ebdd3541c093a0dd1c9d2607ad9de46ee17a6f38600bc2e0100b8fecc72fb024b6102c53a3924e8b7e2c163cfe6e5b84299c4481efa0580694ca8ca0000ba4358b81da547ba46f897dcc83cd6c94a2aff075e25d2f92d420b090ece5728b8fb0000bcf07ee1e924c93e5ce98564b648115a8de881106ad0901cd32cae515986206818f60000be1c8831d018fed582688b5e75be3c9983ac1761fb24a916451c6999c0c66f5e14960000be7d342b6dceaa0bc80f0e1c6d8415a589908ddb50ef59f1ebbd520727ff6079eccc0000c031d8519d06c58a2b42ae7022de4f3242d3938e855501455b3057a28fe38e6c60950000c056d91fd8fb8dabe669748d3e0aebf560449daee47c4e4c1bb6e4cd30d910452c330000c08d5de7a5d97bea2c7ddf516d0635bddc43f326ae2f80e2595b49d4a08c46194cf90000c1f215e6d1c4cede0f67c75dfe8c1f6f3244a318b4405fb775e76b8b6453cd3c10fe0000c22eff63f91830d4f1b882bad30ab87660afb93dba8c027933706a95d063b066086b0000c2ba6d32d980ecb8450511c34f3c4f9a1a46c40cbd61c4d4120de3525d9f057374fe0000c2c65d95e24107e87efcdfd8516890320d0bbb2f2c2c9c79509d74b56d866d61d8c70000c2eea43fd0e45e0e756130e01667533bcaa001e29e0192501e7ce2186ec3554a4cfe0000c5e80accf4092ea6f8ed087544576dddcfdd51366b492868f73b0c9ca19c5f3164ff0000c63599d8442f3faf02474552f5981cff2e70144b5a3bb975b48dec2df1d4062c08980000c65de6003709aa5a6b81354c00fb13e281ac05e852cb4194c69f78566e8ac82854fb0000c6860f7015c86b7074a74feba7e8a19d9fc885521498d7668a12d37efe65480984990000c6bfc6a3cade13b0fc2d1ad1741f3fb4562b89c7f3bfb5622dcdae02db6534555cfd0000c885a011dbc63183416fa77f78a967ad0ffbbbb62e444e9eb490abd2fc67d32610f90000ca386062a9bca5f23c49b71ae24b4005c33bddd54aac0115517728223d441a18cc970000ca38730919998d8ddcf5e729a8a275769b1ae064d1e0f5528db03eede05d5e24a4c40000cabcb84f5dd2bc5f057a9cd2228148bd28b7a4803fd4da26e4595946c744686748670000cc2c72ebc8af12b8ff769b61b56f5add49ff0f4badaa3998ba2b2fb37196786754970000ce4761ea5d8d29dd92a78dc606246d8146320a5e9292de240ab15f60c7802e4a84c60000d22f4b3f7a9f0878a49954e7b4491ca841f68831c2e8aeb383cd79dcdc00295b14c80000d44533a4d21fd9d6f5d57c8cd05c61a6f23f9131cec8ae386b6b437db399ec3d202f0100d45f6bb9da856765d3a631e49424df35f3fbcde96c002e296b989cae415423553cfa0000d49cb030a944c5d66de3ca7d4b8f5b37df03a610a028efdf313b5b405b3614a4c4cc0000d4c9b5341b6040e6f218c5476e5bb87a26fd80e03a0ac65b92d4d3eb917c8f22e8990000d62a2b80ebcda1b2f14d2a903088759ce56482401fb4130cde32775d6d210a6ab4ff0000d6c29a7c39cee45b0e045a94081bc188ef73be2be086d66aefd850fc7eeacc451c2e0100d6f1e7a6aa706b8f6cb64c978ab2379ac9347889e53daa9ad5ac3f62cb24e40428320000d85f93f47f94abb8dfaf8c817af3d4b539847141277479506027f06c54cba440f4ce0000d873e5f608496742d7cadd1eabfb8da118a4e12384e23a694a7b71d14f33366e88680000d8e943700c6625ff869b258c77221ea5fe0a01380faf6059d28be5f69405fd5114eb0000d9ba7812c46be60e86e1dce7cc23721d4cdcaeb25bd8a7960a5fbc48d68c4108b8fb0000da77fa0d82d2d92ee03531ad36f3867ba4fe7ad8bbf7e9bc917c7a767c64ca34fc990000de064386a3512bba067ef6bd099e7f3a3e68fe6074a8c1ab872ddd7beb8c9002ec950000de1b663515515f63ee1228024bf16550f8ecc057c3c11212c9b1988f7261f86bd0fc0000de7d8da7b555d388f637270fe722557129b90be3844921a4b3e01f412a98f54914c80000e040bdc0cc8eaa648c83644a86cdae61f2b2e2b835c5ea77eb5f681c84f456bf9cc70000e043d8f7872cd895f8957c9179c4264816be3e649713cb3bdc523f752602cc3accc90000e23862ea58297bc541940d558dc215aeda8816d2fa9d4245fce260647df60066f09b0000e2b4012c3ee5585e9f326b9d46d05a3e25dbe5e08897d10b26741c0ce0d5053934990000e2d406c4177e89a0ba49b541e1418f7cbccc0b861dd17ceaa73ef7f9882d952058f70000e2f8aef0d1bc20724c28440d40e2cd31e83805a1f641fbb73cc1c5b777346076a4fb0000e40db41d7f07b2b867fae0d7b8ed6767d23f7b53b032710dc5b5043474bf1d119cfe0000e60b065719963027baac7cb4d68145539a0b428aa74048ec1c32b6b25553064608bb0000e693f8c8c6043a5d8c8ed64d56523d157625011947a8a79881987d9e9100963a14cd0000e82d3b7fbff7f5ff1010b279435012e0ce3b68dc387eb2827fc2e9d5001ada7cbcca0000e8382680e672b8403c57f2bd1073c34219fbd40160e8907ff4cbc548976d263faccb0000e8928812930a7d996416662e5c41fd543ba1768d0deb243e17811071e57da3603c640000e8b6a4eefeb4942bbca2bee6baee73280d49c1b7aff8d1aa7616d06bc173ad7fc0990000e8c7ad65c15fa3ba64424a61b177382a0c5468135aecca9ca454f5e7ce4d305b38f90000e8e6242dc60a384a09384160e375c37415e5cec52471dfa32860625b3794e00938f90000eaab0cb55c147ffaf184a4c00513e85f6d5bb6416994fbdd0dd168f3c59a291b30ca0000ecdd548c83457ab43caf7867e2bef91ef783025db9659afd89794ec1220acf29e0970000ee080855f606cce66bdfffb8a73c54a440fa4a4ea1f9a487b7e2dadedaac205beccc0000eebc7887720ec1ce8b759629cb425df5f15011f0d455bfe1a22ad4cfb36d1a3bbc980000f0cab0c194935f449b3baef742e992ebb801fe38ab5c345d7a6195740399cb50ccfb0000f0de782e8bad3c663be60812f0a2ac63464f5da3ec448c73334c07d71ef27f2ce0970000f2e8a164a3517d03fc58abedb1b8c6ed956c691b34ab28f86b3c45e1d8b7777af4c90000f663b33ca66051a034684493aec66df1466d94004a7459ba5de8e9e7daa6c30168650000f6b2256178f0210557ddcc138e8b264204e4674fb2d73a8190faf4ac3724722f00c80000f6c8f87b052f5e857f4bbc467c41c9a60d1a1f2f6b87342ec0159cd8c665084dc0fd0000f85bd4ec9a558cde3f05e33b0b74f9e732cc41f904894873247d4c435a3b8b6300f50000f8720c905d7ac1acab25c4f353df9eb759e0141e4732540d163ac444260f017794ca0000f8beb38300350558332de98cb43eed57db2fec85bd6153c7199ea1416084372e28c80000fa23c23e7352d3aa5fd58289020dd732ad36a009a1c65f3d0132e867b5ee1102c0940000fa393fa644be44f299ba092f4950cc4b358de4c1a05c0a8e673c4f7f6be9c14fbc980000fc6f8380646bfa19f4dc7c1ed6ebfb0a93f5781793aef9224a05e805426d151c98670000fc9ada02baa2ef7a2a24a91f4fdadbbe83a22733dd3ae2deda97f45b4c30dc671c010100fcfd589d8df6da23f65a51de867fac9490ead3ffbb36ce8d1946cec1789a9a46bc2e0100feb586a6d2406aeb67d8dc3895f09dde186e89556ec3c166b6fb1d94bd7a025e542d0100feffabd8f3747b5ca19766866bd59856a2e65c537c6456c6f6863c026915f4732cfb0000"
	palletName := types.Staking
	storageEntry := types.StakingErasRewardPoints
	decodeStorage, err := client.DynamicDecodeStorage(palletName, storageEntry, raw, metaData)
	if err != nil {
		panic(err)
	}
	fmt.Println(decodeStorage)
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
	accountInfo, err := client.SystemAccount("esnGd4DPqtJ7scVJaBsToDJEnYqu9Eo3ib7yE2T4GBPrRnqcw")
	if err != nil {
		panic(err)
	}
	fmt.Println(accountInfo.Nonce)
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
