package Handlers

import (
	model "CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"CitadelCore/AuthorisationServer/SRP"
)

func HandleLogonChallenge(dto model.LogonChallenge,
	repository Repository.AuthorisationRepository) (model.LogonChallengeResponse, *SRP.SRP6) {
	accountinfo := repository.GetAccountInformation(dto.GetAccountName())
	srp, err := SRP.StartSRP(accountinfo.Accountname, accountinfo.Salt[:], accountinfo.Verifier[:])
	generator, prime := SRP.GetParameters()

	if err != nil {
		panic(err)
	}

	response := model.LogonChallengeResponse{}

	response.Command = model.AuthLogonChallenge
	response.ProtocolVersion = 0
	response.Status = model.Success
	var saltarray [32]byte
	copy(saltarray[:], accountinfo.Salt)
	response.Salt = saltarray
	response.Flags = 0
	response.GeneratorSize = 1
	response.Generator = uint8(generator.Int64())
	response.PrimeSize = 32
	var primearray [32]byte
	copy(primearray[:], prime.Bytes()[:32])
	response.Prime = primearray
	var epharray [32]byte
	copy(epharray[:], srp.EphemeralPublicB.Bytes())
	response.EphemeralPublicB = epharray

	return response, srp
}
