package Model

type ReconnectChallengeResponse struct {
	Command   uint8
	Status    uint8
	Salt      [16]byte
	Challenge [16]byte
}
