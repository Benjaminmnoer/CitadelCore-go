package Handlers

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/Legacy"
	"fmt"
)

func HandleReconnectProof(dto Model.ReconnectProof, accountname string) (Model.ReconnectProofResponse, error) {
	authKey, err := Legacy.GetSessionKey(accountname)

	if err != nil {
		return Model.ReconnectProofResponse{}, err
	}

	return Model.ReconnectProofResponse{}, fmt.Errorf("Not implemented")
}
