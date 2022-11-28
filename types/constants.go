package types

const FFISuccess = 1

type HashType string

const (
	Twox64Concat     HashType = "Twox64Concat"
	Blake2_128Concat HashType = "Blake2_128Concat"
)

// Module prefix

type ModuleName string

const (
	System             ModuleName = "System"
	Staking            ModuleName = "Staking"
	Balances           ModuleName = "Balances"
	Utility            ModuleName = "Utility"
	TransactionPayment ModuleName = "TransactionPayment"
)

//

type StorageKey string

const (
	Events  StorageKey = "Events"
	Account StorageKey = "Account"

	ActiveEra               StorageKey = "ActiveEra"
	ErasValidatorReward     StorageKey = "ErasValidatorReward"
	StakingErasRewardPoints StorageKey = "ErasRewardPoints"
	ErasValidatorPrefs      StorageKey = "ErasValidatorPrefs"
	ErasStakersClipped      StorageKey = "ErasStakersClipped"
	Ledger                  StorageKey = "Ledger"
)

type EventID string

const (
	ExtrinsicSuccess EventID = "ExtrinsicSuccess"
	CodeUpdated      EventID = "CodeUpdated"

	Transfer EventID = "Transfer"
	Deposit  EventID = "Deposit"
	Withdraw EventID = "Withdraw"

	Withdrawn EventID = "Withdrawn"
	Slashed   EventID = "Slashed"
	Rewarded  EventID = "Rewarded"
	Bonded    EventID = "Bonded"
	Unbonded  EventID = "Unbonded"
)

func (e EventID) String() string {
	return string(e)
}

type CallId string

const (
	ReBond        CallId = "rebond"
	PayoutStakers CallId = "payout_stakers"
	CallTransfer  CallId = "transfer"
)
