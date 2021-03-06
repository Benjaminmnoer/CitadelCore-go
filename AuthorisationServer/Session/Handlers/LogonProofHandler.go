package Handlers

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/SRP"
	"CitadelCore/Legacy"
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
	for i, j := 0, len(m2arr)-1; i < j; i, j = i+1, j-1 {
		m2arr[i], m2arr[j] = m2arr[j], m2arr[i]
	}
	response.M2 = m2arr
	response.SurveyId = 0
	response.LoginFlags = 0
	var accountflagsarray [4]byte
	hex, _ := hex.DecodeString("00800000")
	copy(accountflagsarray[:], hex)
	for i, j := 0, len(accountflagsarray)-1; i < j; i, j = i+1, j-1 {
		accountflagsarray[i], accountflagsarray[j] = accountflagsarray[j], accountflagsarray[i]
	}
	response.AccountFlags = accountflagsarray

	temp := srp.SessionKey.Bytes()
	for i, j := 0, len(temp)-1; i < j; i, j = i+1, j-1 {
		temp[i], temp[j] = temp[j], temp[i]
	}

	err = Legacy.SetSessionKey(temp, srp.Username)
	if err != nil {
		return Model.LogonProofResponse{}, fmt.Errorf("Couldn't set auth key in legacy database.\n%s\n", err)
	}

	return response, nil
}
