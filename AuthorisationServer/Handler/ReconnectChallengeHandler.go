package Handlers

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	"CitadelCore/AuthorisationServer/Model"
)

var versionChallenge = [16]byte{0xBA, 0xA3, 0x1E, 0x99, 0xA0, 0x0B, 0x21, 0x57, 0xFC, 0x37, 0x3F, 0xB3, 0x69, 0xCD, 0xD2, 0xF1}

func HandleReconnectChallenge(dto Model.LogonChallenge) (Model.ReconnectChallengeResponse, error) {
	bytes := Cryptography.GetRandomBytes(16)
	var salt [16]byte
	copy(salt[:], bytes)
	result := Model.ReconnectChallengeResponse{
		Command:   Model.AuthReconnectChallenge,
		Status:    Model.Success,
		Salt:      salt,
		Challenge: versionChallenge,
	}

	return result, nil
}
