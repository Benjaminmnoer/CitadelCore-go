package SRP_test

import (
	SRP "CitadelCore/AuthorisationServer/srp"
	"encoding/hex"
	"testing"
)

func TestSrpFunctionality(t *testing.T) {
	SRP.InitializaSRP()
	name := "TEST"
	salt, _ := hex.DecodeString("2DC3EC5243A37FC6E9CF0485F0A9DEB1F9901D845F3FA23CC89F3735B7933A1A")
	verifier, _ := hex.DecodeString("7213FA92EE26CAADDE6EBF9B10C275EEA0603AA16362198E967A86A301B712DB")
	ephemeralpublicA, _ := hex.DecodeString("82F8E5D999E51264D8B4646095B0B2A921A262F3396A909AB6939CB1A6DFB89D")
	for i, j := 0, len(ephemeralpublicA)-1; i < j; i, j = i+1, j-1 {
		ephemeralpublicA[i], ephemeralpublicA[j] = ephemeralpublicA[j], ephemeralpublicA[i]
	}

	// Apparently should be reversed?
	m1, _ := hex.DecodeString("e32af4530c1c04d178756720237df932932628b1")

	srp := SRP.NewSrp()
	srp.StartSRP(name, salt, verifier)

	err := srp.VerifyProof(ephemeralpublicA, m1)

	if err != nil {
		t.Fatalf("Wtf happened here: %s\n", err)
	}
}
