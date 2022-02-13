package model

type AccountInformation struct {
	Accountname string
	Salt        []byte
	Verifier    []byte
}
