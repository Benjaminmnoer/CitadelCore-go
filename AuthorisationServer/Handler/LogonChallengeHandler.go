package Handlers

import (
	model "CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"CitadelCore/AuthorisationServer/SRP"
)

func HandleLogonChallenge(dto model.LogonChallenge,
	repository Repository.AuthorisationRepository) (model.LogonChallengeResponse, SRP.SRP6) {
	accountinfo := repository.GetAccountInformation(dto.GetAccountName())
	srp, err := SRP.StartSRP(accountinfo.Accountname, accountinfo.Salt, accountinfo.Verifier)

	if err != nil {
		panic(err)
	}

	response := model.LogonChallengeResponse{}

	response.Command = model.AuthLogonChallenge
	response.ProtocolVersion = 0
	response.Status = model.Success
	response.Salt = accountinfo.Salt
	response.Flags = 0

	return response, srp
}
