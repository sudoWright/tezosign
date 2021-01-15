package models

type Request struct {
	ID         uint64
	Data       string
	Status     string
	ContractID uint64
}

type Sign struct {
	ID        uint64
	RequestID uint64
	Signature string
}
