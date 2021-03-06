package SRP

import (
	"CitadelCore/AuthorisationServer/Cryptography"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func hexFromBigInt(input *big.Int) string {
	bytes := input.Bytes()
	hex := hex.EncodeToString(bytes)
	upper := strings.ToUpper(hex)
	return upper
}

// Calculates server keys and sets the ... variables
func (srp *SRP6) generateServerKeys() error {
	if srp.EphemeralPrivateB.Cmp(bigIntZero) == 0 {
		// srp.EphemeralPrivateB.SetBytes(Cryptography.GetNonce())
		privateb, _ := hex.DecodeString("F1568D79CF6E35A3A44791A12DFC9A09B2FD1B0C90948D29F747E63991E44919")
		srp.EphemeralPrivateB.SetBytes(privateb)
	}

	if srp.EphemeralPrivateB.Cmp(bigIntZero) <= 0 ||
		multiplier.Cmp(bigIntZero) <= 0 ||
		srp.Verifier.Cmp(bigIntZero) <= 0 ||
		Prime.Cmp(bigIntZero) <= 0 ||
		Generator.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("Dont know enough")
	}

	// Calculating public b
	term1 := &big.Int{}
	term2 := &big.Int{}
	term1.Mul(multiplier, srp.Verifier)
	term2.Exp(Generator, srp.EphemeralPrivateB, Prime)
	srp.EphemeralPublicB.Add(term1, term2)
	srp.EphemeralPublicB.Mod(srp.EphemeralPublicB, Prime)

	if srp.Salt.Cmp(bigIntZero) <= 0 ||
		srp.Verifier.Cmp(bigIntZero) <= 0 ||
		srp.EphemeralPrivateB.Cmp(bigIntZero) <= 0 ||
		srp.EphemeralPublicB.Cmp(bigIntZero) <= 0 {

		return fmt.Errorf("Error in setting SRP values.\nSalt: %v\nVerifier: %v\nEphemeral Public B: %v\n", srp.Salt, srp.Verifier, srp.EphemeralPublicB)
	}
	return nil
}

// Constructs the proof that the server knows the key. Sets the M2 variable
func (srp *SRP6) createServerProof() {
	res := Cryptography.Sha1Multiplebytes(getLittleEndian(srp.ephemeralPublicA), getLittleEndian(srp.M1), getLittleEndian(srp.SessionKey))
	setReverseEndian(srp.M2, res)
}

// Verifies the received client proof of key knowledge. Returns error if it does not verify.
func (srp *SRP6) verifyClientProof() error {
	if srp.SessionKey.Cmp(bigIntZero) <= 0 {
		return fmt.Errorf("Trying to prove client without knowing key")
	}

	primtemp := Prime.Bytes()
	for i, j := 0, len(primtemp)-1; i < j; i, j = i+1, j-1 {
		primtemp[i], primtemp[j] = primtemp[j], primtemp[i]
	}

	nHash := Cryptography.Sha1Bytes(primtemp)
	for i, j := 0, len(nHash)-1; i < j; i, j = i+1, j-1 {
		nHash[i], nHash[j] = nHash[j], nHash[i]
	}

	gHash := Cryptography.Sha1Bytes(Generator.Bytes())
	for i, j := 0, len(gHash)-1; i < j; i, j = i+1, j-1 {
		gHash[i], gHash[j] = gHash[j], gHash[i]
	}

	ntemp := &big.Int{}
	gtemp := &big.Int{}
	ntemp.SetBytes(nHash)
	gtemp.SetBytes(gHash)
	dest := make([]byte, 20)
	for i := 0; i < 20; i++ {
		dest[i] = nHash[i] ^ gHash[i]
	}
	for i, j := 0, len(dest)-1; i < j; i, j = i+1, j-1 {
		dest[i], dest[j] = dest[j], dest[i]
	}

	xor := &big.Int{}
	xor.SetBytes(dest)

	uHash := Cryptography.Sha1Bytes([]byte(srp.Username))
	utemp := &big.Int{}
	utemp.SetBytes(uHash)

	m1res := Cryptography.Sha1Multiplebytes(dest, uHash, getLittleEndian(srp.Salt), getLittleEndian(srp.ephemeralPublicA), getLittleEndian(srp.EphemeralPublicB), getLittleEndian(srp.SessionKey))
	temp := &big.Int{}
	setReverseEndian(temp, m1res)
	fmt.Printf("Calculated M1: %s\n", hex.EncodeToString(m1res))

	for i, v := range srp.M1.Bytes() {
		if v != m1res[i] {
			return fmt.Errorf("Could not verify clients proof at index %d!\nExpected %d\nReceived%d\n", i, m1res, srp.M1.Bytes())
		}
	}

	return nil
}

func (srp *SRP6) calculateU() error {
	keyhash := Cryptography.Sha1Multiplebytes(getLittleEndian(srp.ephemeralPublicA), getLittleEndian(srp.EphemeralPublicB))
	setReverseEndian(srp.U, keyhash)

	if srp.U.Cmp(bigIntZero) <= 0 {
		return errors.New("srp: Error setting u value")
	}
	return nil
}

func (srp *SRP6) calculateSessionKey() error {
	temp := big.NewInt(0)

	temp.Exp(srp.Verifier, srp.U, Prime)
	temp.Mul(srp.ephemeralPublicA, temp)
	temp.Exp(temp, srp.EphemeralPrivateB, Prime)
	srp.PreSessionKey.Exp(srp.Verifier, srp.U, Prime)
	srp.PreSessionKey.Mul(srp.ephemeralPublicA, srp.PreSessionKey)
	srp.PreSessionKey.Exp(srp.PreSessionKey, srp.EphemeralPrivateB, Prime)

	// Shamelessly stolen directly from trinitycore.
	sbytes := srp.PreSessionKey.Bytes()
	for i2, j := 0, len(sbytes)-1; i2 < j; i2, j = i2+1, j-1 {
		sbytes[i2], sbytes[j] = sbytes[j], sbytes[i2]
	}

	buffer0 := make([]byte, 16)
	buffer1 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		buffer0[i] = sbytes[2*i+0]
		buffer1[i] = sbytes[2*i+1]
	}

	p := 0
	for ; p < 32; p++ {
		if sbytes[p] != 0 {
			break
		}
	}
	if p%2 == 1 {
		p++
	}
	p /= 2

	hash0 := Cryptography.Sha1Bytes(buffer0[p:])
	hash1 := Cryptography.Sha1Bytes(buffer1[p:])

	res := make([]byte, 40)

	for i := 0; i < 20; i++ {
		res[2*i+0] = hash0[i]
		res[2*i+1] = hash1[i]
	}

	setReverseEndian(srp.SessionKey, res)

	if srp.PreSessionKey.Cmp(bigIntZero) <= 0 || srp.SessionKey.Cmp(bigIntZero) <= 0 {
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

func (srp *SRP6) PrintSRP() {
	fmt.Printf("SRP session for %s\n", srp.Username)
	fmt.Printf("Prime:\nHex %s\nDec %s\n", hexFromBigInt(Prime), Prime.Text(10))
	fmt.Printf("Multiplier:\nHex %s\nDec %s\n", hexFromBigInt(multiplier), multiplier.Text(10))
	fmt.Printf("Generator:\nHex %s\nDec %s\n", hexFromBigInt(Generator), Generator.Text(10))
	fmt.Printf("Verifier:\nHex %s\nDec %s\n", hexFromBigInt(srp.Verifier), srp.Verifier.Text(10))
	fmt.Printf("Salt:\nHex %s\nDec %s\n", hexFromBigInt(srp.Salt), srp.Salt.Text(10))
	fmt.Printf("EphemeralPublicA:\nHex %s\nDec %s\n", hexFromBigInt(srp.ephemeralPublicA), srp.ephemeralPublicA.Text(10))
	fmt.Printf("EphemeralPrivateB:\nHex %s\nDec %s\n", hexFromBigInt(srp.EphemeralPrivateB), srp.EphemeralPrivateB.Text(10))
	fmt.Printf("EphemeralPublicB:\nHex %s\nDec %s\n", hexFromBigInt(srp.EphemeralPublicB), srp.EphemeralPublicB.Text(10))
	fmt.Printf("u:\nHex %s\nDec %s\n", hexFromBigInt(srp.U), srp.U.Text(10))
	fmt.Printf("Pre-Session key:\nHex %s\nDec %s\n", hexFromBigInt(srp.PreSessionKey), srp.PreSessionKey.Text(10))
	fmt.Printf("Session key:\nHex %s\nDec %s\n", hexFromBigInt(srp.SessionKey), srp.SessionKey.Text(10))
	fmt.Printf("M1:\nHex %s\nDec %s\n", hexFromBigInt(srp.M1), srp.M1.Text(10))
	fmt.Printf("M2:\nHex %s\nDec %s\n", hexFromBigInt(srp.M2), srp.M2.Text(10))
	fmt.Println()
}

func getLittleEndian(b *big.Int) []byte {
	temp := b.Bytes()
	for i, j := 0, len(temp)-1; i < j; i, j = i+1, j-1 {
		temp[i], temp[j] = temp[j], temp[i]
	}
	return temp
}

func setReverseEndian(b *big.Int, value []byte) {
	temp := value
	for i, j := 0, len(temp)-1; i < j; i, j = i+1, j-1 {
		temp[i], temp[j] = temp[j], temp[i]
	}
	b.SetBytes(temp)
}
