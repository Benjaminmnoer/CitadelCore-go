package Handlers

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/Legacy"
	"bytes"
	"fmt"
)

func HandleReconnectProof(dto Model.ReconnectProof, accountname string, reconnectproof [16]byte) (Model.ReconnectProofResponse, error) {
	authKey, err := Legacy.GetSessionKey(accountname)

	if err != nil {
		return Model.ReconnectProofResponse{}, err
	}
	hash := Cryptography.Sha1Multiplebytes([]byte(accountname), dto.R1[:], reconnectproof[:], authKey)
	if bytes.Compare(hash, dto.R2[:]) != 0 {
		return Model.ReconnectProofResponse{}, fmt.Errorf("Not equal in value")
	}

	return Model.ReconnectProofResponse{
		Command:    Model.AuthReconnectProof,
		Status:     Model.Success,
		LoginFlags: 0,
	}, nil
}
