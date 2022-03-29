package Model

type ReconnectProof struct {
	Commandd uint8
	R1       [16]byte
	R2       [20]byte
	R3       [20]byte
	NumKeys  uint8
}
