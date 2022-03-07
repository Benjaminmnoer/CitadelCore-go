package Model

type AccountInformation struct {
	Id          uint64
	Accountname string
	Salt        []byte
	Verifier    []byte
}
