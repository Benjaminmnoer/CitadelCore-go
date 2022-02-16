package SRP_test

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	"CitadelCore/AuthorisationServer/SRP"
	"math/big"
	"testing"
)

func TestSrpFunctionality(t *testing.T) {
	name := "TEST"
	salt := []byte{171, 239, 21, 183, 61, 201, 18, 176, 15, 54, 118, 4, 106, 254, 81, 192, 209, 242, 122, 158, 24, 128, 71, 209, 35, 241, 187, 82, 46, 57, 117, 228}
	verifier := []byte{142, 57, 119, 41, 61, 78, 200, 84, 128, 254, 219, 244, 153, 19, 98, 190, 103, 109, 231, 3, 242, 43, 104, 235, 84, 211, 224, 179, 52, 21, 155, 44}

	srp := SRP.NewSrp()
	srp.StartSRP(name, salt, verifier)

	gen, prime := SRP.GetConstants()

	ephemeralprivatea := big.NewInt(0)
	ephemeralprivatea.SetBytes(Cryptography.GetNonce())
	ephemeralpublica := big.NewInt(0)
	ephemeralpublica.Exp(gen, ephemeralprivatea, prime)
}
