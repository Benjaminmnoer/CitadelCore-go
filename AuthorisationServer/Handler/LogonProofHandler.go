package Handlers

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/SRP"
)

func HandleLogonProof(dto Model.LogonProof, srp *SRP.SRP6) Model.LogonProofResponse {
	response := Model.LogonProofResponse{}
	err := srp.VerifyProof(dto.A[:], dto.M1[:])

	if err != nil {
		panic(err)
	}

	response.Command = Model.AuthLogonProof
	response.Error = uint8(Model.Success)
	var m2arr [20]byte
	copy(m2arr[:], srp.M2.Bytes())
	response.M2 = m2arr
	response.SurveyId = 0
	response.LoginFlags = 0
	return response
}
