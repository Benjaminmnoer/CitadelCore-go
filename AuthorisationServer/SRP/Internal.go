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
	privateb, _ := hex.DecodeString("5780F2F000EE71219B4CAD40F268A6AAFB570164AB97CF3AC90725D8C4DD85D0")
	srp.EphemeralPrivateB.SetBytes(privateb)
	// srp.ephemeralPrivateB.SetBytes(Cryptography.GetNonce())

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
func (srp *SRP6) createServerProof() error {
	Cryptography.Sha1Multiplebytes(srp.ephemeralPublicA.Bytes(), srp.M1.Bytes(), srp.PreSessionKey.Bytes())
	return errors.New("Not implemented")
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
	fmt.Printf("NHash:\nHex %s\nDec %s\n", hexFromBigInt(ntemp), ntemp.Text(10))
	fmt.Printf("gHash:\nHex %s\nDec %s\n", hexFromBigInt(gtemp), gtemp.Text(10))
	dest := make([]byte, 20)
	for i := 0; i < 20; i++ {
		dest[i] = nHash[i] ^ gHash[i]
	}
	for i, j := 0, len(dest)-1; i < j; i, j = i+1, j-1 {
		dest[i], dest[j] = dest[j], dest[i]
	}

	xor := &big.Int{}
	xor.SetBytes(dest)
	fmt.Printf("xor:\nHex %s\nDec %s\n", hexFromBigInt(xor), xor.Text(10))
	// xor := make([]byte, Cryptography.Sha1Size())
	// xorlength := safeXORBytes(xor, nHash[:], gHash[:])
	// if xorlength != Cryptography.Sha1Size() {
	// 	return fmt.Errorf("XOR had %d bytes instead of %d", xorlength, Cryptography.Sha1Size())
	// }
	// groupHash := Cryptography.Sha1Bytes(dest)
	// for i, j := 0, len(groupHash)-1; i < j; i, j = i+1, j-1 {
	// 	groupHash[i], groupHash[j] = groupHash[j], groupHash[i]
	// }

	uHash := Cryptography.Sha1Bytes([]byte(srp.Username))
	utemp := &big.Int{}
	utemp.SetBytes(uHash)
	fmt.Printf("Usernamehash:\nHex %s\nDec %s\n", hexFromBigInt(utemp), utemp.Text(10))

	// m1res := Cryptography.Sha1Multiplebytes(groupHash, uHash, srp.Salt.Bytes(), srp.ephemeralPublicA.Bytes(), srp.EphemeralPublicB.Bytes(), srp.SessionKey.Bytes())
	m1res := Cryptography.Sha1Multiplebytes(dest, uHash, getLittleEndian(srp.Salt), getLittleEndian(srp.ephemeralPublicA), getLittleEndian(srp.EphemeralPublicB), getLittleEndian(srp.SessionKey))
	temp := &big.Int{}
	temp.SetBytes(m1res)
	fmt.Printf("Calculated m1:\nHex: %s\nDec: %s\n", hexFromBigInt(temp), temp.Text(10))

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
	fmt.Printf("v ^ u mod N:\nHex %s\nDec %s\n", hexFromBigInt(temp), temp.Text(10))
	temp.Mul(srp.ephemeralPublicA, temp)
	fmt.Printf("a * base:\nHex %s\nDec %s\n", hexFromBigInt(temp), temp.Text(10))
	temp.Exp(temp, srp.EphemeralPrivateB, Prime)
	fmt.Printf("base ^ b mod N:\nHex %s\nDec %s\n", hexFromBigInt(temp), temp.Text(10))
	srp.PreSessionKey.Exp(srp.Verifier, srp.U, Prime)
	srp.PreSessionKey.Mul(srp.ephemeralPublicA, srp.PreSessionKey)
	srp.PreSessionKey.Exp(srp.PreSessionKey, srp.EphemeralPrivateB, Prime)

	// Shamelessly stolen directly from trinitycore
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
	temp0 := &big.Int{}
	temp1 := &big.Int{}
	temp0.SetBytes(buffer0)
	temp1.SetBytes(buffer1)
	fmt.Printf("Buffer0: %s\n", temp0.Text(10))
	fmt.Printf("Buffer1: %s\n", temp1.Text(10))

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
	fmt.Printf("Offset: %d\n", p)

	hash0 := Cryptography.Sha1Bytes(buffer0[p:])
	hash1 := Cryptography.Sha1Bytes(buffer1[p:])

	temp0.SetBytes(hash0)
	temp1.SetBytes(hash1)
	fmt.Printf("Buffer0: %s\n", temp0.Text(10))
	fmt.Printf("Buffer1: %s\n", temp1.Text(10))

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

func (srp *SRP6) printSRP() {
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
