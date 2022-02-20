package SRP

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func (srp *SRP6) isPublicValid() error {
	hexa := hex.EncodeToString(srp.ephemeralPublicA.Bytes())
	hexb := hex.EncodeToString(srp.EphemeralPublicB.Bytes())
	fmt.Printf("A: %s\nB: %s\n", hexa, hexb)
	aresult := big.Int{}
	bresult := big.Int{}

	if aresult.Mod(srp.ephemeralPublicA, Prime); aresult.Sign() == 0 {
		return fmt.Errorf("Public a is not valid")
	}

	if aresult.GCD(nil, nil, srp.ephemeralPublicA, Prime).Cmp(bigIntZero) != 0 {
		return fmt.Errorf("Public a is not valid")
	}

	if bresult.Mod(srp.EphemeralPublicB, Prime); bresult.Sign() == 0 {
		return fmt.Errorf("Public b is not valid")
	}

	if bresult.GCD(nil, nil, srp.EphemeralPublicB, Prime).Cmp(bigIntZero) != 0 {
		return fmt.Errorf("Public b is not valid")
	}

	return nil
}

func hexFromBigInt(input *big.Int) string {
	bytes := input.Bytes()
	hex := hex.EncodeToString(bytes)
	lower := strings.ToUpper(hex)
	// return strings.TrimLeft(lower, "0")
	return lower
}

// Calculates server keys and sets the ... variables
func (srp *SRP6) generateServerKeys() error {
	privateb, _ := hex.DecodeString("254C5F31A0C7622F4D56FAB1B0CB137DA72604717EC72307240BADABCA5EDA81")
	srp.ephemeralPrivateB.SetBytes(privateb)
	// srp.ephemeralPrivateB.SetBytes(Cryptography.GetNonce())

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
	term1.Mul(big.NewInt(3), srp.verifier)
	// term1.Mod(term1, Prime)
	term2.Exp(big.NewInt(7), srp.ephemeralPrivateB, Prime)
	srp.EphemeralPublicB.Add(term1, term2)
	srp.EphemeralPublicB.Mod(srp.EphemeralPublicB, Prime)
	rev := srp.EphemeralPublicB.Bytes()
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	srp.EphemeralPublicB.SetBytes(rev)

	if srp.salt.Cmp(bigIntZero) <= 0 ||
		srp.verifier.Cmp(bigIntZero) <= 0 ||
		srp.ephemeralPrivateB.Cmp(bigIntZero) <= 0 ||
		srp.EphemeralPublicB.Cmp(bigIntZero) <= 0 {

		return fmt.Errorf("Error in setting SRP values.\nSalt: %v\nVerifier: %v\nEphemeral Public B: %v\n", srp.salt, srp.verifier, srp.EphemeralPublicB)
	}
	return nil
}

// Constructs the proof that the server knows the key. Sets the M2 variable
func (srp *SRP6) createServerProof() error {
	Cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.m1.Bytes(), srp.preSessionKey.Bytes())
	return errors.New("Not implemented")
}

// Verifies the received client proof of key knowledge. Returns error if it does not verify.
func (srp *SRP6) verifyClientProof() error {
	if srp.sessionKey.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("Trying to prove client without knowing key")
	}

	// First lets work on the H(H(A) ⊕ H(g)) part.
	nHash := Cryptography.Sha1Bytes(Prime.Bytes())
	gHash := Cryptography.Sha1Bytes(Generator.Bytes())
	xor := make([]byte, Cryptography.Sha1Size())
	xorlength := safeXORBytes(xor, nHash[:], gHash[:])
	if xorlength != Cryptography.Sha1Size() {
		return fmt.Errorf("XOR had %d bytes instead of %d", xorlength, Cryptography.Sha1Size())
	}
	groupHash := Cryptography.Sha1Bytes(xor)

	uHash := Cryptography.Sha1Bytes([]byte(srp.Username))

	m1res := Cryptography.Sha1Multiplebytes(groupHash, uHash, srp.salt.Bytes(), srp.ephemeralPublicA.Bytes(), srp.EphemeralPublicB.Bytes(), srp.sessionKey.Bytes())

	for i, v := range srp.m1.Bytes() {
		if v != m1res[i] {
			return fmt.Errorf("Could not verify clients proof at index %d!\nExpected %d\nReceived%d\n", i, m1res, srp.m1.Bytes())
		}
	}

	return nil
}

func (srp *SRP6) calculateU() error {
	// temp := &big.Int{}
	// temp.Add(srp.ephemeralPublicA, srp.EphemeralPublicB)
	// abytes := srp.ephemeralPublicA.Bytes()
	// for i, j := 0, len(abytes)-1; i < j; i, j = i+1, j-1 {
	// 	abytes[i], abytes[j] = abytes[j], abytes[i]
	// }
	// bbytes := srp.EphemeralPublicB.Bytes()
	// for i, j := 0, len(bbytes)-1; i < j; i, j = i+1, j-1 {
	// 	bbytes[i], bbytes[j] = bbytes[j], bbytes[i]
	// }
	// publica := strings.ToUpper(hex.EncodeToString(abytes))
	// publicb := strings.ToUpper(hex.EncodeToString(bbytes))
	publica := hexFromBigInt(srp.ephemeralPublicA)
	publicb := hexFromBigInt(srp.EphemeralPublicB)
	fmt.Printf("Public A: %s\nPublic B: %s\n", publica, publicb)

	// srp.u.SetBytes(Cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.EphemeralPublicB.Bytes()))
	temp := []byte(fmt.Sprintf("%s%s", publica, publicb))
	// for i, j := 0, len(temp)-1; i < j; i, j = i+1, j-1 {
	// 	temp[i], temp[j] = temp[j], temp[i]
	// }
	hash := Cryptography.Sha1Bytes(temp)
	// for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
	// 	hash[i], hash[j] = hash[j], hash[i]
	// }

	srp.u.SetBytes(hash)
	// h := sha1.New()
	// h.Write([]byte(fmt.Sprintf("%s%s", publica, publicb)))
	// srp.u.SetBytes(Cryptography.Sha1Bytes(temp.Bytes()))

	if srp.u.Cmp(bigIntZero) <= 0 {
		return errors.New("srp: Error setting u value")
	}
	return nil
}

func (srp *SRP6) calculateSessionKey() error {
	base := &big.Int{}

	base.Exp(srp.verifier, srp.u, Prime)
	base.Mul(base, srp.ephemeralPublicA)

	srp.preSessionKey.Exp(base, srp.ephemeralPrivateB, Prime)

	// Stolen directly from trinitycore
	sbytes := srp.preSessionKey.Bytes()
	bufffer0 := make([]byte, 16)
	bufffer1 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		bufffer0[i] = sbytes[2*i+0]
		bufffer1[i] = sbytes[2*i+1]
	}

	i := 0
	for ; i < 32; i++ {
		if sbytes[i] == 0 {
			break
		}
	}
	if i%2 == 1 {
		i++
	}
	i /= 2

	hash0 := Cryptography.Sha1Bytes(bufffer0)
	hash1 := Cryptography.Sha1Bytes(bufffer1)

	res := make([]byte, 40)

	for i := 0; i < 20; i++ {
		res[2*i+0] = hash0[i]
		res[2*i+1] = hash1[i]
	}

	// srp.sessionKey.SetBytes(Cryptography.Sha1Bytes(srp.preSessionKey.Bytes()))
	srp.sessionKey.SetBytes(res)

	if srp.preSessionKey.Cmp(bigIntZero) <= 0 || srp.sessionKey.Cmp(bigIntZero) <= 0 {
		return errors.New("srp: Error setting key values")
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

func (srp *SRP6) printSRP() {
	fmt.Printf("SRP session for %s\n", srp.Username)
	fmt.Printf("Prime: %s\n", hexFromBigInt(Prime))
	fmt.Printf("Multiplier: %s\n", hexFromBigInt(multiplier))
	fmt.Printf("Generator: %s\n", hexFromBigInt(Generator))
	fmt.Printf("Verifier: %s\n", hexFromBigInt(srp.verifier))
	fmt.Printf("Salt: %s\n", hexFromBigInt(srp.salt))
	fmt.Printf("EphemeralPublicA: %s\n", hexFromBigInt(srp.ephemeralPublicA))
	fmt.Printf("EphemeralPrivateB: %s\n", hexFromBigInt(srp.ephemeralPrivateB))
	fmt.Printf("EphemeralPublicB: %s\n", hexFromBigInt(srp.EphemeralPublicB))
	fmt.Printf("u: %s\n", hexFromBigInt(srp.u))
	fmt.Printf("Pre-Session key: %s\n", hexFromBigInt(srp.preSessionKey))
	fmt.Printf("Session key: %s\n", hexFromBigInt(srp.sessionKey))
	fmt.Printf("M1: %s\n", hexFromBigInt(srp.m1))
	fmt.Printf("M2: %s\n", hexFromBigInt(srp.M2))
	fmt.Println()
}
