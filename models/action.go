package models

type ActionType string

const (
	Transfer      ActionType = "transfer"
	Delegation    ActionType = "delegation"
	FATransfer    ActionType = "fa_transfer"
	StorageUpdate ActionType = "storage_update"
	CustomPayload ActionType = "custom"
)
