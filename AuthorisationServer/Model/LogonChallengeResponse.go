package model

type LogonChallengeResponse struct {
	Command          uint8
	ProtocolVersion  uint8
	Status           AuthorisationResult
	EphemeralPublicB []byte
	GeneratorSize    uint8
	Generator        uint8
	PrimeSize        uint8
	Prime            []byte
	Salt             []byte
	CRC              [16]byte
	Flags            uint32 // How many?
}
