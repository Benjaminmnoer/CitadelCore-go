package Model

type AccountInformation struct {
	Accountname string
	Salt        []byte
	Verifier    []byte
}
