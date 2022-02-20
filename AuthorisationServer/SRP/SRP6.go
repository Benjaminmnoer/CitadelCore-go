package SRP

import (
	"encoding/hex"
	"fmt"
	"math/big"
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
	EphemeralPublicB  *big.Int // Must be public or create getter
	preSessionKey     *big.Int
	sessionKey        *big.Int
	u                 *big.Int
	verifier          *big.Int
	salt              *big.Int
	m1                *big.Int
	M2                *big.Int // Must be public or create getter
	Username          string
	badstate          bool
}

func InitializaSRP() {
	res, err := hex.DecodeString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7")

	// Apparently should be reversed? Apparently not?
	// for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
	// 	res[i], res[j] = res[j], res[i]
	// }

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

// Returns pointer to empty srp. Needs to call StartSRP to initialize session.
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
	srp.Username = name
	srp.salt.SetBytes(s)
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
	srp.verifier.SetBytes(v)

	if srp.salt.Cmp(bigIntZero) <= 0 || srp.verifier.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("Error setting salt or verifier")
	}

	error := srp.generateServerKeys()

	if error != nil {
		return error
	}

	srp.printSRP()

	return nil
}

func (srp SRP6) VerifyProof(ephemeralPublicA []byte, m1 []byte) error {
	srp.ephemeralPublicA.SetBytes(ephemeralPublicA)
	srp.m1.SetBytes(m1)

	if srp.ephemeralPublicA.Cmp(bigIntZero) <= 0 || srp.m1.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("srp: Error setting proof values.\nA: %s\nM1: %s\n", hex.EncodeToString(ephemeralPublicA), hex.EncodeToString(m1))
	}

	err := srp.calculateU()
	if err != nil {
		return err
	}

	err = srp.calculateSessionKey()
	if err != nil {
		return err
	}

	err = srp.generateServerKeys()
	if err != nil {
		return err
	}

	// err = srp.verifyClientProof()
	// srp.M2.SetBytes(Cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.m1.Bytes(), srp.sessionKey.Bytes()))

	srp.printSRP()

	return err
}
