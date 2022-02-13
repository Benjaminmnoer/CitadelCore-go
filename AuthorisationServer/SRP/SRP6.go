package SRP

import (
	cryptography "CitadelCore/AuthorisationServer/Cryptography"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

var (
	bigIntZero = big.NewInt(0)
	Generator  = big.NewInt(7)
	Prime      = big.NewInt(0)
	multiplier = big.NewInt(3)
)

type SRP6 struct {
	ephemeralPublicA  *big.Int
	ephemeralPrivateB *big.Int
	ephemeralPublicB  *big.Int
	preSessionKey     *big.Int
	sessionKey        *big.Int
	u                 *big.Int
	verifier          *big.Int
	salt              *big.Int
	m1                *big.Int
	m2                *big.Int
	username          string
}

func InitializaSRP() {
	res, err := hex.DecodeString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7")

	if err != nil {
		fmt.Printf("Error initalizing SRP, %s\n", err)
		return
	}

	Prime = Prime.SetBytes(res)
}

func GetParameters() (generator []byte, prime []byte) {
	return Generator.Bytes(), Prime.Bytes()
}

func StartSRP(name string, s []byte, v []byte) (SRP6, error) {
	srp := SRP6{}

	srp.username = name
	srp.salt.SetBytes(s)
	srp.verifier.SetBytes(s)
	srp.ephemeralPrivateB.SetBytes(cryptography.GetNonce())

	// Calculating public b
	term1 := &big.Int{}
	term2 := &big.Int{}
	term1.Mul(multiplier, srp.verifier)
	term1.Mod(term1, Prime)
	term2.Exp(Generator, srp.ephemeralPrivateB, Prime)
	srp.ephemeralPublicB.Add(term1, term2)

	return srp, nil
}

func (srp SRP6) VerifyProof(ephemeralPublicA []byte, m1 []byte) error {
	srp.ephemeralPublicA.SetBytes(ephemeralPublicA)
	srp.m1.SetBytes(m1)

	if srp.ephemeralPublicA.Cmp(bigIntZero) == 0 || srp.m1.Cmp(bigIntZero) == 0 {
		return errors.New("srp: Error setting proof values")
	}

	temp := &big.Int{}
	temp.Add(srp.ephemeralPublicA, srp.ephemeralPublicB)
	srp.u.SetBytes(cryptography.Sha1Bytes(temp.Bytes()))

	base := &big.Int{}

	base.Exp(srp.verifier, srp.u, Prime)
	base.Mul(base, srp.ephemeralPublicA)

	srp.preSessionKey.Exp(base, srp.ephemeralPrivateB, Prime)
	srp.sessionKey.SetBytes(cryptography.Sha1Bytes(srp.preSessionKey.Bytes()))

	temp = &big.Int{}
	bytes := append(append(srp.ephemeralPublicA.Bytes(), srp.m1.Bytes()...), srp.sessionKey.Bytes()...)
	srp.m2.SetBytes(cryptography.Sha1Bytes(bytes))

	return nil
}
