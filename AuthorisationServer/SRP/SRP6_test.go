package SRP_test

import (
	SRP "CitadelCore/AuthorisationServer/SRP"
	"encoding/hex"
	"strings"
	"testing"
)

func TestSrpFunctionality(t *testing.T) {
	SRP.InitializaSRP()
	name := "TEST"
	salt, _ := hex.DecodeString("ABEF15B73DC912B00F3676046AFE51C0D1F27A9E188047D123F1BB522E3975E4")
	for i, j := 0, len(salt)-1; i < j; i, j = i+1, j-1 {
		salt[i], salt[j] = salt[j], salt[i]
	}

	verifier, _ := hex.DecodeString("2C9B1534B3E0D354EB682BF203E76D67BE621399F4DBFE8054C84E3D2977398E")
	privateb, _ := hex.DecodeString("F1568D79CF6E35A3A44791A12DFC9A09B2FD1B0C90948D29F747E63991E44919")
	ephemeralpublicA, _ := hex.DecodeString("80A902907937B749FF11FD6DA01FFA1446388E3DFC79C2FCFE6BB033D5CC13C1")
	m1, _ := hex.DecodeString(strings.ToUpper("E930033E74341E5B73B27993C092B906820FDF2F"))

	srp := SRP.NewSrp()
	srp.EphemeralPrivateB.SetBytes(privateb)
	srp.StartSRP(name, salt, verifier)

	ephemeralpublicB := strings.ToUpper(hex.EncodeToString(srp.EphemeralPublicB.Bytes()))
	if ephemeralpublicB != "3D9B530FF7CBEE98BEDFE0503BE67FDA9FF0B10C84C9DBF0C83713C0A3917805" {
		t.Fatalf("Wrong public B key.\nExpected: %s\nactual: %s\n", "3D9B530FF7CBEE98BEDFE0503BE67FDA9FF0B10C84C9DBF0C83713C0A3917805", ephemeralpublicB)
	}

	err := srp.VerifyProof(ephemeralpublicA, m1)

	if err != nil {
		t.Fatalf("Wtf happened here: %s\n", err)
	}

	u := strings.ToUpper(hex.EncodeToString(srp.U.Bytes()))
	if u != "A6960BA2DE0C042BABE026A2E42C57161B3791D7" {
		t.Fatalf("Wrong scrambling.\nExpected: %s\nactual: %s\n", "A6960BA2DE0C042BABE026A2E42C57161B3791D7", u)
	}

	presessionkey := strings.ToUpper(hex.EncodeToString(srp.PreSessionKey.Bytes()))
	if presessionkey != "82D4F223C030BB430CDA9213F35751FFE3541F31446DD4516CCB4AE0069D1D80" {
		t.Fatalf("Wrong pre session key.\nExpected: %s\nactual: %s\n", "82D4F223C030BB430CDA9213F35751FFE3541F31446DD4516CCB4AE0069D1D80", presessionkey)
	}

	sessionkey := strings.ToUpper(hex.EncodeToString(srp.SessionKey.Bytes()))
	if sessionkey != "7C4E350B088D742221DECAD55E28582A9AA3752F32765DC8AB3467C92B3A63008F519F48A59B9862" {
		t.Fatalf("Wrong session key.\nExpected: %s\nactual: %s\n", "7C4E350B088D742221DECAD55E28582A9AA3752F32765DC8AB3467C92B3A63008F519F48A59B9862", sessionkey)
	}

	m2 := strings.ToUpper(hex.EncodeToString(srp.M2.Bytes()))
	if m2 != "CEC637390B518A728A9264AE52A6144F58B5E383" {
		t.Fatalf("Wrong m2.\nExpected: %s\nactual: %s\n", "CEC637390B518A728A9264AE52A6144F58B5E383", m2)
	}
}
