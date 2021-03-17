package models

type ActionType string

const (
	Transfer      ActionType = "transfer"
	Delegation    ActionType = "delegation"
	FATransfer    ActionType = "fa_transfer"
	FA2Transfer   ActionType = "fa2_transfer"
	StorageUpdate ActionType = "storage_update"
	CustomPayload ActionType = "custom"

	//Vesting
	VestingVest        ActionType = "vesting_vest"
	VestingSetDelegate ActionType = "vesting_set_delegate"

	//Income transfer
	IncomeTransfer ActionType = "income_transfer"
)
