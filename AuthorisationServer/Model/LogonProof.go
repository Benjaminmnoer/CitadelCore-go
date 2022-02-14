package Model

type LogonProof struct {
	Command   uint8
	A         [32]byte
	M1        [20]byte
	CRC       [20]byte
	NKeys     uint8
	TwoFactor uint8
}
