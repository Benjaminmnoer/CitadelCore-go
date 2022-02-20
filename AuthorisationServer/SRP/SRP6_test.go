package SRP_test

import (
	"CitadelCore/AuthorisationServer/SRP"
	"encoding/hex"
	"testing"
)

func TestSrpFunctionality(t *testing.T) {
	SRP.InitializaSRP()
	name := "TEST"
	salt, _ := hex.DecodeString("abef15b73dc912b00f3676046afe51c0d1f27a9e188047d123f1bb522e3975e4")
	verifier, _ := hex.DecodeString("8e3977293d4ec85480fedbf4991362be676de703f22b68eb54d3e0b334159b2c")
	ephemeralpublicA, _ := hex.DecodeString("1B482488FFD5CEFE78C1510530D77F671A0E712635ABA1069455E2B6CC76692A")
	// Apparently should be reversed?
	m1, _ := hex.DecodeString("e32af4530c1c04d178756720237df932932628b1")

	srp := SRP.NewSrp()
	srp.StartSRP(name, salt, verifier)

	err := srp.VerifyProof(ephemeralpublicA, m1)

	if err != nil {
		t.Fatalf("Wtf happened here: %s\n", err)
	}
}
