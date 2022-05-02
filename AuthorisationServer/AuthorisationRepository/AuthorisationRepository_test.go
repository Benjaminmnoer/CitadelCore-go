package AuthorisationRepository_test

import (
	"CitadelCore/AuthorisationServer/AuthorisationRepository"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGetAccountInformation(t *testing.T) {
	accountname := "TEST"
	repo := AuthorisationRepository.InitializeAuthorisationRepository()

	result := repo.GetAccountInformation(accountname)

	resacc := result.Accountname

	fmt.Println(resacc)
	fmt.Println(hex.EncodeToString(result.Salt))
	fmt.Println(hex.EncodeToString(result.Verifier))
	if resacc != accountname {
		t.Fatalf("Account name did not return as expected. Expected %s, was %s\n", accountname, resacc)
	}
}

func TestGetRealm(t *testing.T) {
	repo := AuthorisationRepository.InitializeAuthorisationRepository()

	result, _ := repo.GetRealms()

	fmt.Println(result)
	for _, v := range result {
		fmt.Println(v.Name)
	}
}
