package Model

type LogonChallengeResponse struct {
	Command          uint8
	ProtocolVersion  uint8
	Status           AuthorisationResult
	EphemeralPublicB [32]byte
	GeneratorSize    uint8
	Generator        uint8
	PrimeSize        uint8
	Prime            [32]byte
	Salt             [32]byte
	CRC              [16]byte
	Flags            uint32 // How many?
}
