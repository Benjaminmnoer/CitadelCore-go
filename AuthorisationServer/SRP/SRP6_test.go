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
	ephemeralpublicA, _ := hex.DecodeString("01DAED2C5CB059CB7B503E9E2C6CB1CBA39EC009E786048BDF838663485820C0")
	// for i, j := 0, len(ephemeralpublicA)-1; i < j; i, j = i+1, j-1 {
	// 	ephemeralpublicA[i], ephemeralpublicA[j] = ephemeralpublicA[j], ephemeralpublicA[i]
	// }
	m1, _ := hex.DecodeString(strings.ToUpper("8F021980098667C9476FE8B188001629ECD120E6"))

	srp := SRP.NewSrp()
	srp.StartSRP(name, salt, verifier)

	// ephemeralpublicB := strings.ToUpper(hex.EncodeToString(srp.EphemeralPublicB.Bytes()))
	// if ephemeralpublicB != "5045EAB01C37CF7427FF834BC5E6351B3198B1E4F52E16BC0A0C51AAE4F30BB6" {
	// 	t.Fatalf("Wrong public B key.\nExpected: %s\nactual: %s\n", "5045EAB01C37CF7427FF834BC5E6351B3198B1E4F52E16BC0A0C51AAE4F30BB6", ephemeralpublicB)
	// }

	srp.VerifyProof(ephemeralpublicA, m1)

	// if err != nil {
	// 	t.Fatalf("Wtf happened here: %s\n", err)
	// }

	// u := strings.ToUpper(hex.EncodeToString(srp.U.Bytes()))
	// if u != "4D2759CB07F303C949BA58F814CE5D979FBF2CE7" {
	// 	t.Fatalf("Wrong scrambling.\nExpected: %s\nactual: %s\n", "4D2759CB07F303C949BA58F814CE5D979FBF2CE7", u)
	// }

	// presessionkey := strings.ToUpper(hex.EncodeToString(srp.PreSessionKey.Bytes()))
	// if presessionkey != "87CBAFA0B2371A2C6A95717112679B9731DAFCD70E439AEC3AF8CF7DBB02DEF4" {
	// 	t.Fatalf("Wrong pre session key.\nExpected: %s\nactual: %s\n", "87CBAFA0B2371A2C6A95717112679B9731DAFCD70E439AEC3AF8CF7DBB02DEF4", presessionkey)
	// }

	// sessionkey := strings.ToUpper(hex.EncodeToString(srp.SessionKey.Bytes()))
	// if sessionkey != "50055EACE70C08BC649927B294166F58B4B295461E738F843E35DBBDFD4C3FB2CAB0B805F2ED5003" {
	// 	t.Fatalf("Wrong session key.\nExpected: %s\nactual: %s\n", "50055EACE70C08BC649927B294166F58B4B295461E738F843E35DBBDFD4C3FB2CAB0B805F2ED5003", sessionkey)
	// }

	// m2 := strings.ToLower(hex.EncodeToString(srp.M2.Bytes()))
	// if u != "7A1CEE70435D89D0D7E2B6D86D590B5EF680E986" {
	// 	t.Fatalf("Wrong m2.\nExpected: %s\nactual: %s\n", "7A1CEE70435D89D0D7E2B6D86D590B5EF680E986", m2)
	// }
}
