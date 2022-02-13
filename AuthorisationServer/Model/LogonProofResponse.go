package model

type LogonProofResponse struct {
	Command      uint8
	Error        uint8
	M2           [20]byte
	AccountFlags uint32
	SurveyId     uint32
	LoginFlags   uint16
}
