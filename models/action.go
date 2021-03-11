package models

type ActionType string

const (
	Transfer      ActionType = "transfer"
	Delegation    ActionType = "delegation"
	FATransfer    ActionType = "fa_transfer"
	FA2Transfer   ActionType = "fa2_transfer"
	StorageUpdate ActionType = "storage_update"
	CustomPayload ActionType = "custom"

	//Income transfer
	IncomeTransfer ActionType = "income_transfer"
)
