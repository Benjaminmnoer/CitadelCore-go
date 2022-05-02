package Handlers_test

import (
	"CitadelCore/AuthorisationServer/AuthorisationRepository"
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/SRP"
	"CitadelCore/AuthorisationServer/Session/Handlers"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestLoginHandlers(t *testing.T) {
	// TODO: Use mock or test db instead.
	repo := AuthorisationRepository.InitializeAuthorisationRepository()
		SRP.InitializaSRP()
	srp := SRP.NewSrp()
	privateb, _ := hex.DecodeString("1679C11B7751DF582C7560BF13EC0AF56720A25DA754F1A4C9018DD959E5E3AA")
	srp.EphemeralPrivateB.SetBytes(privateb)
	anbytes := *new([32]byte)
	copy(anbytes[:], []byte("TEST"))
	challenge := Model.LogonChallenge{
		Accountnamelength: 4,
		Accountname:       anbytes,
	}

	fmt.Printf("Testing for name: %s\n", challenge.GetAccountName())

	challengeresponse := Handlers.HandleLogonChallenge(challenge, repo, srp)

	proof := Model.LogonProof{}
	abytes := *new([32]byte)
	ahex, _ := hex.DecodeString("ebf58553e50098726fb10d7dd782c67a026512e3f1a66f3ef6b44a6f37c30e7e")
	copy(abytes[:], ahex)
	proof.A = abytes
	m1bytes := *new([20]byte)
	m1hex, _ := hex.DecodeString("cfb802a1f20fcff47736fd54e859b3da273afb86")
	copy(m1bytes[:], m1hex)
	proof.M1 = m1bytes
	proofresponse, err := Handlers.HandleLogonProof(proof, srp)
	srp.PrintSRP()
	if err != nil {
		t.Fatalf("Failed with error: %s\n", err)
	}

	if hex.EncodeToString(challengeresponse.EphemeralPublicB[:]) != "bd360a8594e65f763fb95d29118b39aca6d2cd8b6fd3f14f4c66fe97a69c224c" {
		t.Fatalf("Wrong public b value returned.")
	}

	if hex.EncodeToString(challengeresponse.Prime[:]) != "b79b3e2a87823cab8f5ebfbf8eb10108535006298b5badbd5b53e1895e644b89" {
		t.Fatalf("Wrong prime value returned.")
	}

	if hex.EncodeToString(challengeresponse.Salt[:]) != "abef15b73dc912b00f3676046afe51c0d1f27a9e188047d123f1bb522e3975e4" {
		t.Fatalf("Wrong salt value returned.")
	}

	if hex.EncodeToString(proofresponse.M2[:]) != "feaa247f1f525d7a957c02cf1973757c4210ebf1" {
		t.Fatalf("Wrong m2 value returned.")
	}
}
