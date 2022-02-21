package SRP

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	"crypto/sha1"
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
	return lower
	// return strings.TrimLeft(lower, "0")
}

// Calculates server keys and sets the ... variables
func (srp *SRP6) generateServerKeys() error {
	privateb, _ := hex.DecodeString("469E8A23A75A0C3C26D9E7BFB8CBAA8B8654C3DA1C24111D8B24989AB6984021")
	// setLittleEndian(srp.ephemeralPrivateB, privateb)
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
	term1.Mul(multiplier, srp.verifier)
	// term1.Mod(term1, Prime)
	term2.Exp(Generator, srp.ephemeralPrivateB, Prime)
	srp.EphemeralPublicB.Add(term1, term2)
	srp.EphemeralPublicB.Mod(srp.EphemeralPublicB, Prime)
	temp := srp.EphemeralPublicB.Bytes()
	srp.EphemeralPublicB.SetBytes(temp)

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
	// srp.u.SetBytes(Cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.EphemeralPublicB.Bytes()))
	// publica := strings.ToLower(hexFromBigInt(srp.ephemeralPublicA))
	// publicb := strings.ToLower(hexFromBigInt(srp.EphemeralPublicB))
	tempa := srp.ephemeralPublicA.Bytes()
	tempb := srp.EphemeralPublicB.Bytes()
	hasher := sha1.New()
	hasher.Write(append(tempa, tempb...))
	// hasher.Write([]byte(fmt.Sprintf("%s%s", publica, publicb)))
	// hasher.Write(srp.ephemeralPublicA.Bytes())
	// hasher.Write(srp.EphemeralPublicB.Bytes())
	res := hasher.Sum(nil)
	srp.u.SetBytes(res)

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

func getLittleEndian(b *big.Int) []byte {
	temp := b.Bytes()
	for i, j := 0, len(temp)-1; i < j; i, j = i+1, j-1 {
		temp[i], temp[j] = temp[j], temp[i]
	}
	return temp
}

func setLittleEndian(b *big.Int, value []byte) {
	for i, j := 0, len(value)-1; i < j; i, j = i+1, j-1 {
		value[i], value[j] = value[j], value[i]
	}
	b.SetBytes(value)
}
