package Model_test

import (
	model "CitadelCore/AuthorisationServer/Model"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGetAccountName(t *testing.T) {
	name := "TEST"
	input, _ := hex.DecodeString("54455354")
	var arr [32]byte
	copy(arr[:], input[:])
	logonchallenge := model.LogonChallenge{Accountnamelength: 4, Accountname: arr}

	result := logonchallenge.GetAccountName()

	fmt.Println(result)

	if name != result {
		t.Fatalf("Expected %s, got %s\n", name, result)
	}
}
