package Cryptography

import "crypto/sha1"

func Sha1String(input string) []byte {
	sha := sha1.Sum([]byte(input))
	return sha[:]
}

func Sha1Bytes(input []byte) []byte {
	sha := sha1.Sum(input)
	return sha[:]
}

func Sha1Multiplebytes(inputs ...[]byte) []byte {
	sha := sha1.New()
	for _, v := range inputs {
		sha.Write(v)
	}
	return sha.Sum(nil)
}

func Sha1Size() int {
	return sha1.Size
}
