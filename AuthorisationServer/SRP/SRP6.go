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
	EphemeralPrivateB *big.Int
	EphemeralPublicB  *big.Int // Must be public or create getter
	PreSessionKey     *big.Int
	SessionKey        *big.Int
	U                 *big.Int
	Verifier          *big.Int
	Salt              *big.Int
	M1                *big.Int
	M2                *big.Int // Must be public or create getter
	Username          string
	Badstate          bool
}

func InitializaSRP() {
	res, err := hex.DecodeString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7")

	if err != nil {
		fmt.Printf("Error initalizing SRP, %s\n", err)
		return
	}

	Prime = Prime.SetBytes(res)
}

// Returns pointer to empty srp. Needs to call StartSRP to initialize session.
func NewSrp() *SRP6 {
	return &SRP6{
		ephemeralPublicA:  big.NewInt(0),
		EphemeralPrivateB: big.NewInt(0), // Public for testing purposes.
		EphemeralPublicB:  big.NewInt(0),
		PreSessionKey:     big.NewInt(0),
		SessionKey:        big.NewInt(0),
		U:                 big.NewInt(0),
		Verifier:          big.NewInt(0),
		Salt:              big.NewInt(0),
		M1:                big.NewInt(0),
		M2:                big.NewInt(0),
		Username:          "",
		Badstate:          false,
	}
}

// Starts a SRP session for a single user.
func (srp *SRP6) StartSRP(name string, s []byte, v []byte) error {
	srp.Username = name
	srp.Salt.SetBytes(s)
	srp.Verifier.SetBytes(v)

	if srp.Salt.Cmp(bigIntZero) <= 0 || srp.Verifier.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("SRP: Error setting salt or verifier.")
	}

	err := srp.generateServerKeys()

	if err != nil {
		return err
	}

	return nil
}

func (srp SRP6) VerifyProof(ephemeralPublicA []byte, m1 []byte) error {
	srp.ephemeralPublicA.SetBytes(ephemeralPublicA)
	srp.M1.SetBytes(m1)

	if srp.ephemeralPublicA.Cmp(bigIntZero) <= 0 || srp.M1.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("SRP: Error setting logon proof values.\nA: %s\nM1: %s\n", hex.EncodeToString(ephemeralPublicA), hex.EncodeToString(m1))
	}

	err := srp.calculateU()
	if err != nil {
		return err
	}

	err = srp.calculateSessionKey()
	if err != nil {
		return err
	}

	err = srp.verifyClientProof()
	if err != nil {
		return err
	}

	srp.createServerProof()

	return nil
}
