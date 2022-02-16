package Model

type LogonProofResponse struct {
	Command      uint8
	Error        uint8
	M2           [20]byte
	AccountFlags [4]byte
	SurveyId     uint32
	LoginFlags   uint16
}
