package Handlers

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/SRP"
	"encoding/hex"
)

func HandleLogonChallenge(dto Model.LogonChallenge, srp *SRP.SRP6) Model.LogonChallengeResponse {
	accountinfo := repository.GetAccountInformation(dto.GetAccountName())
	err := srp.StartSRP(accountinfo.Accountname, accountinfo.Salt[:], accountinfo.Verifier[:])

	generator := SRP.Generator
	prime := SRP.Prime

	if err != nil {
		panic(err)
	}

	response := Model.LogonChallengeResponse{}

	response.Command = Model.AuthLogonChallenge
	response.ProtocolVersion = 0
	response.Status = Model.Success
	var saltarray [32]byte
	copy(saltarray[:], accountinfo.Salt)
	for i, j := 0, len(saltarray)-1; i < j; i, j = i+1, j-1 {
		saltarray[i], saltarray[j] = saltarray[j], saltarray[i]
	}
	response.Salt = saltarray
	response.Flags = 0
	response.GeneratorSize = 1
	response.Generator = uint8(generator.Int64())
	response.PrimeSize = 32
	var primearray [32]byte
	copy(primearray[:], prime.Bytes()[:32])
	for i, j := 0, len(primearray)-1; i < j; i, j = i+1, j-1 {
		primearray[i], primearray[j] = primearray[j], primearray[i]
	}
	response.Prime = primearray
	var epharray [32]byte
	copy(epharray[:], srp.EphemeralPublicB.Bytes())
	for i, j := 0, len(epharray)-1; i < j; i, j = i+1, j-1 {
		epharray[i], epharray[j] = epharray[j], epharray[i]
	}
	response.EphemeralPublicB = epharray
	var crcarray [16]byte
	crchash, err := hex.DecodeString("baa31e99a00b2157fc373fb369cdd2f1") // Hardcoded values is always the best
	copy(crcarray[:], crchash)
	response.CRC = crcarray

	return response
}
