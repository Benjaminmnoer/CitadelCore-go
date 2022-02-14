package Model

type AccountInformation struct {
	Accountname string
	Salt        [32]byte
	Verifier    [32]byte
}
