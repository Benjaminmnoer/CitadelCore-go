package Cryptography

import (
	"fmt"
	"math/rand"
)

const (
	noncesize = 32
)

func GetNonce() []byte {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		fmt.Printf("Error generating nonce, %s\n", err)
		return nil
	}

	return b
}

func GetRandomBytes(i int) []byte {
	b := make([]byte, i)
	_, err := rand.Read(b)

	if err != nil {
		fmt.Printf("Error generating bytes, %s\n", err)
		return nil
	}

	return b
}
