package SRP

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	cryptography "CitadelCore/AuthorisationServer/Cryptography"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	bigIntZero = big.NewInt(0)
	bigIntOne  = big.NewInt(1)
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
	badstate          bool
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
func GetConstants() (*big.Int, *big.Int) {
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
		badstate:          false,
	}
}

// Starts a SRP session for a single user.
func (srp *SRP6) StartSRP(name string, s []byte, v []byte) error {
	// fmt.Printf("Starting SRP session for %s\nVerifier: %d\nSalt: %d\n", name, string(v), string(s))

	srp.Username = name
	srp.salt.SetBytes(s)
	srp.verifier.SetBytes(s)
	error := srp.generateServerKeys()

	if error != nil {
		return error
	}

	fmt.Printf("SRP pointer returned: %p\n", &srp)

	return error
}

func (srp SRP6) VerifyProof(ephemeralPublicA []byte, m1 []byte) error {
	srp.ephemeralPublicA.SetBytes(ephemeralPublicA)
	srp.m1.SetBytes(m1)
	fmt.Printf("A: %s\nM1: %s\n", hexFromBigInt(srp.ephemeralPublicA), hexFromBigInt(srp.m1))

	if srp.ephemeralPublicA.Cmp(bigIntZero) == 0 || srp.m1.Cmp(bigIntZero) == 0 {
		return errors.New("srp: Error setting proof values")
	}

	// err := srp.IsPublicValid()
	// if err != nil {
	// 	return err
	// }

	// temp := &big.Int{}
	// temp.Add(srp.ephemeralPublicA, srp.EphemeralPublicB)
	publica := hexFromBigInt(srp.ephemeralPublicA)
	publicb := hexFromBigInt(srp.EphemeralPublicB)
	// srp.u.SetBytes(cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.EphemeralPublicB.Bytes()))
	srp.u.SetBytes(cryptography.Sha1Bytes([]byte(fmt.Sprintf("%s%s", publica, publicb))))

	base := &big.Int{}

	base.Exp(srp.verifier, srp.u, Prime)
	base.Mul(base, srp.ephemeralPublicA)

	srp.preSessionKey.Exp(base, srp.ephemeralPrivateB, Prime)
	srp.sessionKey.SetBytes(cryptography.Sha1Bytes(srp.preSessionKey.Bytes()))

	err := srp.verifyClientProof()
	if err != nil {
		return err
	}

	srp.M2.SetBytes(cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.m1.Bytes(), srp.sessionKey.Bytes()))

	return nil
}

// func (srp *SRP6) IsPublicValid() error {
// 	hexa := hex.EncodeToString(srp.ephemeralPublicA.Bytes())
// 	hexb := hex.EncodeToString(srp.EphemeralPublicB.Bytes())
// 	fmt.Printf("A: %s\nB: %s\n", hexa, hexb)
// 	aresult := big.Int{}
// 	bresult := big.Int{}

// 	if aresult.Mod(srp.ephemeralPublicA, Prime); aresult.Sign() == 0 {
// 		return fmt.Errorf("Public a is not valid")
// 	}

// 	if aresult.GCD(nil, nil, srp.ephemeralPublicA, Prime).Cmp(bigIntZero) != 0 {
// 		return fmt.Errorf("Public a is not valid")
// 	}

// 	if bresult.Mod(srp.EphemeralPublicB, Prime); bresult.Sign() == 0 {
// 		return fmt.Errorf("Public b is not valid")
// 	}

// 	if bresult.GCD(nil, nil, srp.EphemeralPublicB, Prime).Cmp(bigIntZero) != 0 {
// 		return fmt.Errorf("Public b is not valid")
// 	}

// 	return nil
// }

func hexFromBigInt(input *big.Int) string {
	bytes := input.Bytes()
	hex := hex.EncodeToString(bytes)
	lower := strings.ToLower(hex)
	return strings.TrimLeft(lower, "0")
}

// Calculates server keys and sets the ... variables
func (srp *SRP6) generateServerKeys() error {
	srp.ephemeralPrivateB.SetBytes(cryptography.GetNonce())

	if srp.ephemeralPrivateB.Cmp(bigIntZero) <= 0 ||
		multiplier.Cmp(bigIntZero) <= 0 ||
		srp.verifier.Cmp(bigIntZero) <= 0 ||
		Prime.Cmp(bigIntZero) <= 0 ||
		Generator.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("Dont know enough")
	}

	// Calculating public b
	term1 := &big.Int{}
	term2 := &big.Int{}
	term1.Mul(multiplier, srp.verifier)
	term1.Mod(term1, Prime)
	term2.Exp(Generator, srp.ephemeralPrivateB, Prime)
	srp.EphemeralPublicB.Add(term1, term2)
	srp.EphemeralPublicB.Mod(srp.EphemeralPublicB, Prime)

	if srp.salt.Cmp(bigIntZero) <= 0 ||
		srp.verifier.Cmp(bigIntZero) <= 0 ||
		srp.ephemeralPrivateB.Cmp(bigIntZero) <= 0 ||
		srp.EphemeralPublicB.Cmp(bigIntZero) <= 0 {

		return fmt.Errorf("Error in setting SRP values.\nSalt: %v\nVerifier: %v\nEphemeral Public B: %v\n", srp.salt, srp.verifier, srp.EphemeralPublicB)
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
