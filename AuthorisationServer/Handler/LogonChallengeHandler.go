package Handlers

import (
	model "CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"CitadelCore/AuthorisationServer/SRP"
	"encoding/hex"
	"fmt"
)

func HandleLogonChallenge(dto model.LogonChallenge,
	repository Repository.AuthorisationRepository, srp *SRP.SRP6) model.LogonChallengeResponse {
	accountinfo := repository.GetAccountInformation(dto.GetAccountName())
	err := srp.StartSRP(accountinfo.Accountname, accountinfo.Salt[:], accountinfo.Verifier[:])

	fmt.Printf("SRP pointer returned: %p\n", srp)

	generator, prime := SRP.GetConstants()

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
	var crcarray [16]byte
	crchash, err := hex.DecodeString("baa31e99a00b2157fc373fb369cdd2f1")
	copy(crcarray[:], crchash)
	response.CRC = crcarray

	return response
}
