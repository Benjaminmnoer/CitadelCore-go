package Handlers

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/SRP"
	"encoding/hex"
	"fmt"
)

func HandleLogonProof(dto Model.LogonProof, srp *SRP.SRP6) (Model.LogonProofResponse, error) {
	response := Model.LogonProofResponse{}
	a := dto.A[:]
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}

	m1 := dto.M1[:]
	for i, j := 0, len(m1)-1; i < j; i, j = i+1, j-1 {
		m1[i], m1[j] = m1[j], m1[i]
	}

	err := srp.VerifyProof(a, m1)

	if err != nil {
		return response, fmt.Errorf("SRP: Error happened while verifying proof\n%s\n", err.Error())
	}

	response.Command = Model.AuthLogonProof
	response.Error = uint8(Model.Success)
	var m2arr [20]byte
	copy(m2arr[:], srp.M2.Bytes())
	response.M2 = m2arr
	response.SurveyId = 0
	response.LoginFlags = 0
	var accountflagsarray [4]byte
	hex, _ := hex.DecodeString("00080000")
	copy(accountflagsarray[:], hex)
	response.AccountFlags = accountflagsarray

	return response, nil
}
