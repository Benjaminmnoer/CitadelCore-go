package SRP

import (
	"CitadelCore/AuthorisationServer/Cryptography"
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
	EphemeralPublicB  *big.Int // Must be public
	preSessionKey     *big.Int
	sessionKey        *big.Int
	u                 *big.Int
	verifier          *big.Int
	salt              *big.Int
	m1                *big.Int
	M2                *big.Int // Must be public
	Username          string
}

func InitializaSRP() {
	res, err := hex.DecodeString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7")

	// Apparently should be reversed?
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	if err != nil {
		fmt.Printf("Error initalizing SRP, %s\n", err)
		return
	}

	Prime = Prime.SetBytes(res)
}

// Returns server SRP constants.
func GetConstants() (generator *big.Int, prime *big.Int) {
	return Generator, Prime
}

// Returns empty srp. Needs to call StartSRP to initialize session.
func NewSrp() *SRP6 {
	return &SRP6{
		ephemeralPublicA:  big.NewInt(0),
		ephemeralPrivateB: big.NewInt(0),
		EphemeralPublicB:  big.NewInt(0),
		preSessionKey:     big.NewInt(0),
		sessionKey:        big.NewInt(0),
		u:                 big.NewInt(0),
		verifier:          big.NewInt(0),
		salt:              big.NewInt(0),
		m1:                big.NewInt(0),
		M2:                big.NewInt(0),
		Username:          "",
	}
}

// Starts a SRP session for a single user.
func (srp *SRP6) StartSRP(name string, s []byte, v []byte) error {
	// fmt.Printf("Starting SRP session for %s\nVerifier: %d\nSalt: %d\n", name, string(v), string(s))

	srp.Username = name
	srp.salt.SetBytes(s)
	srp.verifier.SetBytes(s)
	error := srp.generatePrivateKeys()

	fmt.Printf("SRP pointer returned: %p\n", &srp)

	return error
}

func (srp SRP6) VerifyProof(ephemeralPublicA []byte, m1 []byte) error {
	srp.ephemeralPublicA.SetBytes(ephemeralPublicA)
	srp.m1.SetBytes(m1)

	if srp.ephemeralPublicA.Cmp(bigIntZero) == 0 || srp.m1.Cmp(bigIntZero) == 0 {
		return errors.New("srp: Error setting proof values")
	}

	temp := &big.Int{}
	temp.Add(srp.ephemeralPublicA, srp.EphemeralPublicB)
	srp.u.SetBytes(cryptography.Sha1Bytes(temp.Bytes()))

	base := &big.Int{}

	base.Exp(srp.verifier, srp.u, Prime)
	base.Mul(base, srp.ephemeralPublicA)

	srp.preSessionKey.Exp(base, srp.ephemeralPrivateB, Prime)
	srp.sessionKey.SetBytes(cryptography.Sha1Bytes(srp.preSessionKey.Bytes()))

	err := srp.verifyClientProof()
	if err != nil {
		return err
	}

	temp = &big.Int{}
	bytes := append(append(srp.ephemeralPublicA.Bytes(), srp.m1.Bytes()...), srp.sessionKey.Bytes()...)
	srp.M2.SetBytes(cryptography.Sha1Bytes(bytes))

	return nil
}

// Calculates server keys and sets the ... variables
func (srp *SRP6) generatePrivateKeys() error {
	srp.ephemeralPrivateB.SetBytes(cryptography.GetNonce())

	// Calculating public b
	term1 := &big.Int{}
	term2 := &big.Int{}
	term1.Mul(multiplier, srp.verifier)
	term1.Mod(term1, Prime)
	term2.Exp(Generator, srp.ephemeralPrivateB, Prime)
	srp.EphemeralPublicB.Add(term1, term2)

	if srp.salt.Cmp(bigIntZero) <= 0 ||
		srp.verifier.Cmp(bigIntZero) <= 0 ||
		srp.ephemeralPrivateB.Cmp(bigIntZero) <= 0 ||
		srp.EphemeralPublicB.Cmp(bigIntZero) <= 0 {
		fmt.Printf("Error in setting SRP values.\nSalt: %v\nVerifier: %v\nEphemeral Public B: %v\n", srp.salt, srp.verifier, srp.EphemeralPublicB)
		return errors.New("Error in setting SRP values")
	}
	return nil
}

func (srp *SRP6) generateCommonKey() error {
	return nil
}

// Constructs the proof that the server knows the key. Sets the M2 variable
// func (srp *SRP6) createServerProof() error {
// 	Cryptography.Sha1Multiplebytes(a, m1, srp.preSessionKey.Bytes())
// 	return errors.New("Not implemented")
// }

// Verifies the received client proof of key knowledge. Returns error if it does not verify.
func (srp *SRP6) verifyClientProof() error {
	if srp.sessionKey.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("Trying to prove client without knowing key")
	}

	// First lets work on the H(H(A) ⊕ H(g)) part.
	nHash := Cryptography.Sha1Bytes(Prime.Bytes())
	gHash := Cryptography.Sha1Bytes(Generator.Bytes())
	xor := make([]byte, cryptography.Sha1Size())
	xorlength := safeXORBytes(xor, nHash[:], gHash[:])
	if xorlength != cryptography.Sha1Size() {
		return fmt.Errorf("XOR had %d bytes instead of %d", xorlength, cryptography.Sha1Size())
	}
	groupHash := cryptography.Sha1Bytes(xor)

	uHash := cryptography.Sha1Bytes([]byte(srp.Username))

	m1res := cryptography.Sha1Multiplebytes(groupHash, uHash, srp.salt.Bytes(), srp.ephemeralPublicA.Bytes(), srp.EphemeralPublicB.Bytes(), srp.sessionKey.Bytes())

	for i, v := range srp.m1.Bytes() {
		if v != m1res[i] {
			return fmt.Errorf("Could not verify clients proof at index %d!\nExpected %d\nReceived%d\n", i, m1res, srp.m1.Bytes())
		}
	}

	return nil
}

func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}
